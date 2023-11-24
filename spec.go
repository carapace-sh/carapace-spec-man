package spec

import (
	"fmt"
	"os/exec"
	"regexp"
	"sort"
	"strings"

	"github.com/lmorg/murex/utils/man"
	"github.com/neurosnap/sentences/english"
	"github.com/rsteube/carapace-spec/pkg/command"
)

func getPages(exe string) (map[string]string, error) {
	output, err := exec.Command("man", "-k", fmt.Sprintf("^%v($|-)", exe)).Output()
	if err != nil {
		return nil, err
	}

	r := regexp.MustCompile(`^(?P<name>.*) \(\d+\) +- (?P<description>.*)$`)

	pages := make(map[string]string)
	for _, line := range strings.Split(string(output), "\n") {
		if matches := r.FindStringSubmatch(line); matches != nil {
			pages[matches[1]] = matches[2]
		}
	}
	return pages, nil
}

type flag struct {
	definition  string
	description string
}

func flatten(m map[string]string) map[string]flag {
	byDescription := make(map[string][]string)
	for key, value := range m {
		byDescription[value] = append(byDescription[value], key)
	}

	flattened := make(map[string]flag)
	for key, value := range byDescription {
		sort.Slice(value, func(i, j int) bool {
			return len(value[i]) < len(value[j])
		})
		name := strings.TrimLeft(value[len(value)-1], "-")
		flattened[name] = flag{
			definition:  strings.Join(value, ", "),
			description: key,
		}
	}
	return flattened
}

func parse(manpage string, trimDescriptions bool) (*command.Command, error) {
	_, m := man.ParseByStdio(strings.NewReader(manpage))

	tokenizer, err := english.NewSentenceTokenizer(nil)
	if err != nil {
		return nil, err
	}

	cmd := &command.Command{
		Flags: make(map[string]string),
	}
	cmd.Completion.Flag = make(map[string][]string)
	cmd.Completion.PositionalAny = []string{"$files"}

	for name, flag := range flatten(m) {
		if strings.HasPrefix(flag.description, "eg: ") {
			if trimDescriptions {
				_, flag.description, _ = strings.Cut(flag.description, "--") // TODO check if found
			}
			flag.definition += "=" // TODO check test after eg: for optarg (`=..`, `[=..]`)

			cmd.Completion.Flag[name] = []string{"$files"}
		}

		if trimDescriptions {
			if tokens := tokenizer.Tokenize(flag.description); len(tokens) > 0 {
				flag.description = tokens[0].Text
				flag.description = strings.TrimSuffix(flag.description, ".")
				flag.description = strings.TrimSpace(flag.description)
			}
		}
		cmd.Flags[flag.definition] = flag.description
	}

	return cmd, nil
}

func loadPage(page string) (string, error) {
	output, err := exec.Command("man", page).Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

type pageMap map[string]string

func (p pageMap) traverse(page string, trimDescriptions bool) (*command.Command, error) {
	content, err := loadPage(page)
	if err != nil {
		return nil, err
	}

	cmd, err := parse(content, trimDescriptions)
	if err != nil {
		return nil, err
	}

	splitted := strings.Split(page, "-")
	cmd.Name = splitted[len(splitted)-1]
	cmd.Description = p[page]
	delete(p, page)

	for name := range p {
		if !strings.HasPrefix(name, page+"-") {
			continue
		}
		if trimmed := strings.TrimPrefix(name, page+"-"); !strings.Contains(trimmed, "-") {
			subcmd, err := p.traverse(name, trimDescriptions)
			if err != nil {
				return nil, err
			}
			cmd.Commands = append(cmd.Commands, *subcmd)
			delete(p, name)
		}
	}
	return cmd, nil
}

func Command(exe string, trimDescriptions bool) (*command.Command, error) {
	pages, err := getPages(exe)
	if err != nil {
		return nil, err
	}
	return pageMap(pages).traverse(exe, trimDescriptions)
}
