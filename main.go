package main

import (
	"fmt"
	"math/rand"
	"os"
	"slices"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func main() {
	fmt.Println("Choosing a today's happy boy...")
	if os.Getenv("DISCORD_TOKEN") == "" {
		fatal("Set the DISCORD_TOKEN environment variable to your bot token.")
	}
	discord, err := discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))
	if err != nil {
		fatal("Error creating Discord session: ", err)
	}

	guilds, err := discord.UserGuilds(100, "", "")
	if err != nil {
		fatal("Error getting guilds: ", err)
	}
	for _, guild := range guilds {
		fmt.Println("Dicing for guild: ", guild.Name, "(", guild.ID, ")")

		fmt.Println("Getting roles")
		roles, err := discord.GuildRoles(guild.ID)
		if err != nil {
			fmt.Println("Error getting roles: ", err)
			continue
		}

		role := findRole(roles)
		if role == nil {
			fmt.Println("No happyboy role found in guild: ", guild.Name)
			continue
		}

		fmt.Println("Getting members")
		members, err := discord.GuildMembers(guild.ID, "", 1000)
		if err != nil {
			fmt.Println("Error getting members: ", err)
			continue
		}

		fmt.Println("Getting channels")
		channels, err := discord.GuildChannels(guild.ID)
		if err != nil {
			fmt.Println("Error getting channels: ", err)
			continue
		}

		// Remove role from all members
		for _, member := range members {
			if !slices.Contains(member.Roles, role.ID) {
				continue
			}
			err := discord.GuildMemberRoleRemove(guild.ID, member.User.ID, role.ID)
			if err != nil {
				fmt.Println("Error removing role: ", err)
				continue
			}
		}

		happyboy := dice(members)

		fmt.Println("Setting happyboy role for: ", happyboy.User.Username)

		err = discord.GuildMemberRoleAdd(guild.ID, happyboy.User.ID, role.ID)
		if err != nil {
			fmt.Println("Error adding role: ", err)
			continue
		}

		channel := findChannel(channels)
		if channel == nil {
			fmt.Println("No happyboy channel found in guild: ", guild.Name)
			continue
		}

		fmt.Println("Sending message to channel: ", channel.Name)
		_, err = discord.ChannelMessageSend(channel.ID, fmt.Sprintf("<@%s>\nYou are Today's Happy Boy!", happyboy.User.ID))
		if err != nil {
			fmt.Println("Error sending message: ", err)
			continue
		}
	}
}

func findChannel(channels []*discordgo.Channel) *discordgo.Channel {
	for _, channel := range channels {
		if isContainHappyBoy(channel.Name) {
			return channel
		}
	}
	return nil
}

func findRole(roles []*discordgo.Role) *discordgo.Role {
	for _, role := range roles {
		if isContainHappyBoy(role.Name) {
			return role
		}
	}
	return nil
}

func isContainHappyBoy(name string) bool {
	lower := strings.ToLower(name)
	if strings.Contains(lower, "bot") {
		return false
	}
	if strings.Contains(lower, "happyboy") {
		return true
	}
	if strings.Contains(lower, "happy boy") {
		return true
	}
	if strings.Contains(name, "ハッピーボーイ") {
		return true
	}
	return false
}

func dice(members []*discordgo.Member) *discordgo.Member {
	n := len(members)
	for {
		member := members[rand.Intn(n)]
		if member.User.Bot {
			continue
		}
		return member
	}
}

func fatal(e ...any) {
	fmt.Println(e...)
	os.Exit(1)
}
