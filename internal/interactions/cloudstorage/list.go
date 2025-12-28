package cloudstorage

import (
	"fmt"

	"github.com/8h9x/BlastBot/internal/sessions"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/google/uuid"
)

func listHandler(event *handler.CommandEvent) error {
	discordId := event.User().ID

	session, err := sessions.GetSessionForUser(discordId)
	if err != nil {
		return fmt.Errorf("unable to get session for user (%s): %s", discordId, err)
	}

	cloudstorageList, err := listUserCloudstorageSorted(session)
	if err != nil {
		return err
	}

	var components []discord.ContainerSubComponent

	for _, cloudstoragePointer := range cloudstorageList {
		components = append(components, discord.NewTextDisplayf("**Name:** %s\n**Uploaded At:** <t:%d:f>\n**Size:** %s\n",
			cloudstoragePointer.Filename,
			cloudstoragePointer.Uploaded.Unix(),
			humanizeBytes(cloudstoragePointer.Length)),
		)
		// components = append(components, discord.NewActionRow(discord.NewPrimaryButton("Download", fmt.Sprintf("/cloudstorage/download/%s", cloudstoragePointer.Filename))))
	}

	err = event.CreateMessage(discord.MessageCreate{
		Flags: discord.MessageFlagIsComponentsV2,
		Components: []discord.LayoutComponent{
			discord.NewContainer(components...),
		},
	})

	return err
}

func isUUID(s string) bool {
	if len(s) >= 36 {
		_, err := uuid.Parse(s[:36])
		return err == nil
	}

	_, err := uuid.Parse(s)
	return err == nil
}

func humanizeBytes(b int) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}

	div, exp := int(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	units := []string{"KB", "MB", "GB", "TB", "PB", "EB"}
	return fmt.Sprintf("%.2f %s", float64(b)/float64(div), units[exp])
}
