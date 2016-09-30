package main

import (
	"math/rand"
)

var (
	_quoteQuips = []string{
		"\"Part of the beauty of me is that I am very rich.\"\n - Donald Trump",
		"\"I don’t wear a ‘rug’—it’s mine. And I promise not to talk about your massive plastic surgeries that didn’t work.\"\n - Donald Trump",
		"\"You know, wealthy people don’t like me.\"\n - Donald Trump",
		"\"When I think I’m right, nothing bothers me.\"\n - Donald Trump",
		"\"I could be happy living in a studio apartment.\"\n - Donald Trump",
		"\"I have a great relationship with the blacks.\"\n - Donald Trump",
		"\"I believe [the media] like making me out to be something more sinister than I really am.\"\n - Donald Trump",
		"\"I will be so good at the military your head will spin.\"\n - Donald Trump",
		"\"I don’t even consider myself ambitious.\"\n - Donald Trump",
		"\"As everybody knows, but the haters & losers refuse to acknowledge, I do not wear a “wig.” My hair may not be perfect but it’s mine.\"\n - Donald Trump",
		"\"Sorry losers and haters, but my I.Q. is one of the highest -and you all know it! Please don’t feel so stupid or insecure, it’s not your fault.\"\n - Donald Trump",
		"\"I put my name on buildings because it sells better. I don’t do it because, gee, I need that.\"\n - Donald Trump",
		"\"I think I’m like the largest or one of the largest [political] contributors. I’m maxed out every year.\"\n - Donald Trump",
		"\"My net worth fluctuates, and it goes up and down with markets and with attitudes and with feelings—even my own feelings—but I try.\"\n - Donald Trump",
		"\"Money is a little bit of a scorecard, but I don’t do it for the money. I do it because I really enjoy it. I love the creative process.\"\n - Donald Trump",
		"\"There have been many bad things said about me over the years, and in some cases they’ve been true. It doesn’t bother me. If I have a fault and somebody exposes that fault or talks about that fault, you won’t hear me complain.\"\n - Donald Trump",
		"\"It’s a great building. It’s the second-tallest building in Chicago, and I always say it was better for the people of Chicago than it was for Donald Trump.\"\n - Donald Trump",
		"\"At Trump Tower, I know everybody that goes up, and everybody that comes down.\"\n - Donald Trump",
		"\"I love Oprah. Oprah would always be my first choice [for Vice President].\"\n - Donald Trump",
		"\"There is something on that birth certificate — maybe religion, maybe it says he’s a Muslim, I don’t know. Maybe he doesn’t want that. Or, he may not have one.\"\n - Donald Trump",
		"\"If I ever ran for office, I’d do better as a Democrat than as a Republican–and that’s not because I’d be more liberal, because I’m conservative. But the working guy would elect me. He likes me. When I walk down the street, those cabbies start yelling out their windows.\"\n - Donald Trump",
		"\"I think if this country gets any kinder or gentler, it’s literally going to cease to exist\"\n - Donald Trump",
		"\"I will be the greatest jobs president God ever created.\"\n - Donald Trump",
		"\"Bing bing, bong bong, bing bing bing.\"\n - Donald Trump",
		"\"Romney — I have a Gucci store that’s worth more than Romney.\"\n - Donald Trump",
		"\"How stupid are the people of Iowa?\"\n - Donald Trump",
		"\"I don’t want to use the word ‘screwed’, but I screwed him.\"\n - Donald Trump",
		"\"I have tremendous respect for the Japanese people, I mean, you can respect somebody that’s beating the hell out of you.\"\n - Donald Trump",
		"\"I’ll shake hands. I shake hands with people. But it’s not something I like — look, I’m not a huge fan of Japan, but I love their custom.\"\n - Donald Trump",
		"\"The concept of global warming was created by and for the Chinese in order to make U.S. manufacturing non-competitive.\"\n - Donald Trump",
		"\"I would never buy Ivana any decent jewels or pictures. Why give her negotiable assets?\"\n - Donald Trump",
		"\"My life has been about winning. My life has not been about losing.\"\n - Donald Trump",
		"\"She does have a very nice figure. I’ve said that if Ivanka weren’t my daughter, perhaps I’d be dating her.\"\n - Donald Trump",
		"\"I want five children, like in my own family, because with five, then I will know that one will be guaranteed to turn out like me.\"\n - Donald Trump",
		"\"If you have the money, having children is great.\"\n - Donald Trump",
		"\"I have never seen a thin person drinking Diet Coke.\"\n - Donald Trump",
		"\"The concept of shaking hands is absolutely terrible, and statistically I’ve been proven right.\"\n - Donald Trump",
		"\"Do you mind if I sit back a little? Because your breath is very bad—it really is.\"\n - Donald Trump",
		"\"In the second grade I actually gave a teacher a black eye — I punched my music teacher because I didn’t think he knew anything about music and I almost got expelled.\"\n - Donald Trump",
		"\"You may get AIDS by kissing.\"\n - Donald Trump",
	}
)

// return a random quip
func randomQuip() string {
	return _quoteQuips[rand.Intn(len(_quoteQuips))]
}
