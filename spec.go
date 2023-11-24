package spec

import (
	"fmt"
	"os/exec"
	"regexp"
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

func parse(manpage string) (*command.Command, error) {
	_, m := man.ParseByStdio(strings.NewReader(manpage))

	tokenizer, err := english.NewSentenceTokenizer(nil)
	if err != nil {
		return nil, err
	}

	cmd := &command.Command{
		Flags: make(map[string]string),
	}
	cmd.Completion.PositionalAny = []string{"$files"}

	for flag, description := range m {
		if strings.HasPrefix(description, "eg: ") {
			_, description, _ = strings.Cut(description, "--") // TODO check if found
		}

		if tokens := tokenizer.Tokenize(description); len(tokens) > 0 {
			description = tokens[0].Text
			description = strings.TrimSuffix(description, ".")
			description = strings.TrimSpace(description)
			cmd.Flags[flag] = description
		}
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

func (p pageMap) traverse(page string) (*command.Command, error) {
	content, err := loadPage(page)
	if err != nil {
		return nil, err
	}

	cmd, err := parse(content)
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
			subcmd, err := p.traverse(name)
			if err != nil {
				return nil, err
			}
			cmd.Commands = append(cmd.Commands, *subcmd)
			delete(p, name)
		}
	}
	return cmd, nil
}

func Command(exe string) (*command.Command, error) {
	pages, err := getPages(exe)
	if err != nil {
		return nil, err
	}
	return pageMap(pages).traverse(exe)
}
