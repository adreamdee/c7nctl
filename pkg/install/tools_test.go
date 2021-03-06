package install

import (
	"fmt"
	"github.com/choerodon/c7nctl/pkg/kube"
	"github.com/vinkdong/gox/log"
	"testing"
)

func TestGetNewsData(t *testing.T) {
	ctx := Context{
		Client:    kube.GetClient(),
		Namespace: "test",
	}
	log.Info(ctx.GetOrCreateConfigMapData(staticLogName, staticLogKey))
}

func TestSaveNewsData(t *testing.T) {
	ctx := Context{
		Client:    kube.GetClient(),
		Namespace: "test",
	}

	news := &News{
		Name:      "testnews2",
		Namespace: "test",
		Type:      PvcType,
		Status:    FailedStatus,
		Reason:    "reason1 ",
	}
	ctx.SaveNews(news)
}

func TestRandomToken(t *testing.T) {
	fmt.Println(RandomToken(17), RandomToken(12))
}
