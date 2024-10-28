package accesstoken_test

import (
	"context"
	"testing"
	"time"

	"github.com/fgrosse/zaptest"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/prometheus/client_golang/prometheus"
	"github.snooguts.net/reddit/achilles-sdk/pkg/fsm/metrics"
	"github.snooguts.net/reddit/achilles-sdk/pkg/io"
	"github.snooguts.net/reddit/achilles-sdk/pkg/logging"
	achratelimiter "github.snooguts.net/reddit/achilles-sdk/pkg/ratelimiter"
	sdktest "github.snooguts.net/reddit/achilles-sdk/pkg/test"
	"github.snooguts.net/reddit/achilles-token-controller/internal/controllers/accesstoken"
	"github.snooguts.net/reddit/achilles-token-controller/internal/controlplane"
	intscheme "github.snooguts.net/reddit/achilles-token-controller/internal/scheme"
	"github.snooguts.net/reddit/achilles-token-controller/internal/test"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	ctrllog "sigs.k8s.io/controller-runtime/pkg/log"
	ctrlzap "sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

var (
	ctx     context.Context
	testEnv *sdktest.TestEnv
	c       client.Client
	scheme  *runtime.Scheme
	log     *zap.SugaredLogger
)

func TestAccessToken(t *testing.T) {
	RegisterFailHandler(Fail)
	ctrllog.SetLogger(ctrlzap.New(ctrlzap.WriteTo(GinkgoWriter), ctrlzap.UseDevMode(true)))
	RunSpecs(t, "AccessToken Suite")
}

var _ = BeforeSuite(func() {
	SetDefaultEventuallyTimeout(15 * time.Second)
	SetDefaultEventuallyPollingInterval(200 * time.Millisecond)

	log = zaptest.LoggerWriter(GinkgoWriter).Sugar()
	ctx = logging.NewContext(context.Background(), log)
	rl := achratelimiter.NewDefaultProviderRateLimiter(achratelimiter.DefaultProviderRPS)

	scheme = intscheme.MustNewScheme()

	var err error
	testEnv, err = sdktest.NewEnvTestBuilder(ctx).
		WithCRDDirectoryPaths(
			test.CRDPaths(),
		).
		WithScheme(scheme).
		WithLog(log.Desugar()).
		WithManagerSetupFns(
			func(mgr manager.Manager) error {
				// setup controller being tested
				clientApplicator := &io.ClientApplicator{
					Client:     mgr.GetClient(),
					Applicator: io.NewAPIPatchingApplicator(mgr.GetClient()),
				}

				cpCtx := controlplane.Context{
					Metrics: metrics.MustMakeMetrics(scheme, prometheus.NewRegistry()),
				}

				return accesstoken.SetupController(ctx, cpCtx, mgr, rl, clientApplicator)
			},
		).
		WithKubeConfigFile("./").
		Start()

	Expect(err).ToNot(HaveOccurred())

	c = testEnv.Client
})

var _ = AfterSuite(func() {
	By("tearing down the test environment")
	err := testEnv.Stop()
	Expect(err).NotTo(HaveOccurred())
})
