package verify

import (
	"fmt"
	"time"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"github.com/theckman/yacspin"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/cakehappens/seaworthy/pkg/clioptions"
	"github.com/cakehappens/seaworthy/pkg/kubernetes"
	"github.com/cakehappens/seaworthy/pkg/kubernetes/health"
)

// New creates a new Verify cobra command
func New(streams clioptions.IOStreams, resourcer kubernetes.Resourcer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "verify (-f FILENAME | TYPE [NAME])",
		Short: "verify",
		Run: func(cmd *cobra.Command, args []string) {
			resourceChan := make(chan []unstructured.Unstructured, 1)
			errChan := make(chan error, 1)

			go func() {
				defer func() {
					close(resourceChan)
					close(errChan)
				}()
				resources, err := resourcer(cmd.Context(), func(opt *kubernetes.ResourcerOptions) {
					opt.Type = args[0]
				})

				resourceChan <- resources
				errChan <- err
			}()

			cfg := yacspin.Config{
				Frequency:         100 * time.Millisecond,
				Writer:            streams.ErrOut,
				CharSet:           yacspin.CharSets[14],
				Message:           "Getting data from k8s",
				StopMessage:       "Data retrieved from K8s",
				StopCharacter:     "✓ ",
				StopColors:        []string{"fgGreen"},
				StopFailMessage:   "Failed getting data from K9s",
				StopFailCharacter: "X ",
				StopFailColors:    []string{"fgRed"},
			}

			spinner, err := yacspin.New(cfg)
			if err != nil {
				panic(err)
			}

			_ = spinner.Start()

			err = <-errChan
			resources := <-resourceChan

			if err != nil {
				_ = spinner.StopFail()
				return
			}

			_ = spinner.Stop()

			const resultMessageFormat = "%s %s: %s - %s\n"

			for _, r := range resources {
				status := health.ResourceHealth(r)

				switch code := status.Code; code {
				case health.Healthy:
					_, _ = fmt.Fprintf(streams.Out, resultMessageFormat, color.Green.Sprint("✓"), r.GetName(), code, status.Message)
				case health.Progressing:
					_, _ = fmt.Fprintf(streams.Out, resultMessageFormat, "🔄️ ", r.GetName(), code, status.Message)
				case health.Unsupported:
					_, _ = fmt.Fprintf(streams.Out, resultMessageFormat, "⚠️ ", r.GetName(), code, status.Message)
				case health.Unknown:
					_, _ = fmt.Fprintf(streams.Out, resultMessageFormat, "?️ ", r.GetName(), code, status.Message)
				case health.Degraded:
					_, _ = fmt.Fprintf(streams.Out, resultMessageFormat, "🔻️ ", r.GetName(), code, status.Message)
				case health.Suspended:
					_, _ = fmt.Fprintf(streams.Out, resultMessageFormat, "⏸️ ", r.GetName(), code, status.Message)
				case health.Missing:
					_, _ = fmt.Fprintf(streams.Out, resultMessageFormat, "👤️ ", r.GetName(), code, status.Message)
				default:
					panic(fmt.Errorf("unknown status code: %s", code))
				}
			}
		},
	}

	return cmd
}
