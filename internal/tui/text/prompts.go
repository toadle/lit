package text

import "math/rand"

var queryPrompts = []string{
	"Hit me...",
	"Letters incoming...",
	"Type something...",
	"Got any ideas...",
	"Surprise me...",
	"Let's see...",
	"Let's go...",
	"Your query...",
}

func GetRandomQueryPrompt() string {
	return queryPrompts[rand.Intn(len(queryPrompts))]
}
