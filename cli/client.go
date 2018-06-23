package cli

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/manifoldco/promptui"
	"github.com/pkg/errors"

	"github.com/RobertGrantEllis/t9/client"
)

func runClientWithArgs(args ...string) {

	client := getClientFromArgs(args...)

	exactPrompt := &promptui.Select{
		Label: `Would you like to see all matches starting with the digits you specify, or only exact matches?`,
		Items: []string{`prefix`, `exact`},
	}

	selected, _, err := exactPrompt.Run()
	if err == promptui.ErrEOF || err == promptui.ErrInterrupt {
		return
	} else if err != nil {
		printError(err)
		return
	}

	exact := selected == 1
	if exact {
		fmt.Println(`only exact matches will be returned`)
	} else {
		fmt.Println(`all matches starting with the specified digits will be returned`)
	}

	digitsPrompt := &promptui.Prompt{
		Label:    `Enter Digits`,
		Validate: validateDigits,
	}

	for {

		digits, err := digitsPrompt.Run()
		if err == promptui.ErrEOF || err == promptui.ErrInterrupt {
			return
		} else if err != nil {
			printError(err)
			continue
		}

		digits = strings.TrimSpace(digits)

		start := time.Now()
		words, err := client.SimpleLookup(digits, exact)
		responseTime := time.Now().Sub(start)
		if err != nil {
			fail(err)
		} else {
			fmt.Println(words)
			fmt.Printf("(response time: %s)\n", responseTime)
		}
	}
}

func getClientFromArgs(args ...string) client.Client {

	configuration := getClientConfigurationFromArgs(args...)

	client, err := client.New(*configuration)
	if err != nil {
		fail(err)
	}

	return client
}

var digitsValidRegexp = regexp.MustCompile(`^[2-9]*$`)

func validateDigits(digits string) error {

	digits = strings.TrimSpace(digits)

	if !digitsValidRegexp.MatchString(digits) {
		return errors.New(`enter only digits 2-9`)
	} else if len(digits) < 2 {
		return errors.New(`enter at least 2 digits`)
	}

	return nil
}
