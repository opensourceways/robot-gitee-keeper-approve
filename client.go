package main

type ghclient struct {
	cli iClient
}

func (c *ghclient) BotName() (string, error) {
	bot, err := c.cli.GetBot()
	if err != nil {
		return "", err
	}
	return bot.Login, nil
}

func (c *ghclient) AddLabel(org, repo string, number int32, label string) error {
	return c.cli.AddPRLabel(org, repo, number, label)
}

func (c *ghclient) RemoveLabel(org, repo string, number int, label string) error {
	return c.cli.RemovePRLabel(org, repo, int32(number), label)
}
