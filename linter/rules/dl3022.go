package rules

import (
	"fmt"
	"strings"

	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

// validateDL3022 COPY --from should reference a previously defined FROM alias
func validateDL3022(node *parser.Node, file string) (rst []string, err error) {
	fromImage := ""
	isAs, isAsBuild := false, false
	for _, child := range node.Children {
		if child.Value == FROM {
			for _, v := range strings.Fields(child.Original) {
				switch v {
				case "as":
					isAs = true
				case "build":
					if isAs {
						isAsBuild = true
					}
				default:
					if fromImage == "" && v != "FROM" && v != "from" {
						fromImage = v
					} else if fromImage == v && !isAsBuild {
						rst = append(rst, fmt.Sprintf("%s:%v DL3022 COPY --from should reference a previously defined FROM alias\n", file, child.StartLine))
					}
				}
			}
		}
	}
	return rst, nil
}
