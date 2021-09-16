package deploy

import (
	"fmt"

	"github.com/lovego/cmd"
	"github.com/lovego/xiaomei/release"
)

func rmDeploy(svcName, env, feature string) error {
	script := fmt.Sprintf(`
for name in $(docker ps -af name='%s' --format '{{.Names}}'); do
	docker stop $name >/dev/null 2>&1 && docker rm $name
done
`, release.ContainerNameRegexp(svcName, env))
	return eachNodeRun(env, script, feature)
}

func eachNodeRun(env, script, feature string) error {
	for _, node := range release.GetCluster(env).GetNodes(feature) {
		if _, err := node.Run(cmd.O{}, script); err != nil {
			return err
		}
	}
	return nil
}
