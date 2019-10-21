package words

import (
	"regexp"
)

// GetCensorRegex returns regex of words to censor
func GetCensorRegex() *regexp.Regexp {
	return regexp.MustCompile(`(?i)\bboner\b|child ?p(?:or|ro)+n|cock(?:block|suck|juice|slut)|\bcocks?\b|\bclit|\bcum\b|\bcunts?|\bdeep ?throat\b|\benema\b|\bfag(?:got)?\b|\bh(?:o|0)+rny\b|\bincest|\blabia\b|puss(?:(?:a)?y|ies)\b|\b-rape\b|\brape(?:d|s)?\b|\b(?:(?:butt|ass|gang)? ?rape|\brapi(?:ng|st))\b|gang ?bang|pa?edophil|bea?stiality|(?:blow|hand|rim)+ ?job|j(?:ac|er)+k(?:ing)? ?off|(?:shit|glory) ?hole|(?:suck(?:s|ed|ing)?|lick(?:s|ed|ing)?|strok(?:es?|ed|ing)|thrust(?:s|ed|ing)?|insert(?:s|ed|ing)?)+ (?:his|my|your)? ?(?:big|thick|throbbing|swollen|engorged|erect|hard(?:ened|ening)?)? ?(?:dick|member)|\bnigger|\bsperms?\b|\bsemen\b|\btits\b|adult ?friendfinder|videos ?xxx|\bnud(?:e|ity)\b|\btwat\b|x-?rated|\bturn(?:ing|ed)? me on\b|\breach(?:ed|ing)? (?:his|her)? ?climax\b|\b(?:lick(?:ed)?|finger) (?:my|his|her|your) (?:ass|butt|tight )hole\b|\bthrust(?:s|ed|ing)? (?:himself|in and out)\b|\ban(?:al|us)+\b|\bdicks?\b|\borg(?:y|ies)\b|\bscat\b|\b(?:3|three)some\b|\bwhor(?:es?|ing)\b|\bboob|\bbreast|finger(?:ed|ing)|f(?:-| |\.)u(?:-| |\.)c(?:-| |\.)k|s(?:-| |\.)e(?:-| |\.)x|s(?:-| |\.)m(?:-| |\.)u(?:-| |\.)t|fuu+ck|ve+gina|daddy ?kink|(?:jack|jerk) (?:me|you) off|(?:throbbing|swollen|erect|flaccid|hardened|engorged) (?:member|dick|penis)|\b(?:his|her) slit\b|\bgrinds (?:herself|himself)\b|bukkake|cuckold|cum(?:slut|shake|dump)|cunnilingus|dildo|ejaculat|erection|fellatio|foreplay|fornicat|genital|hentai|horniness|intercourse|jizz|masturbat|orgasm|pre-cum|cumming|strapon|testicles|titties|titjob|creampie|penis|vagina|viagra|dry hump|bdsm|kinky|erotic|prostitut|naked|pervert|porn|coitus|copulat|fetish|nympho|slut|striptease|vibrator|shotacon|陰茎|変態|underage sex|sex|cameltoe|upskirt|jailbait|squirter|futanari|facesitting|fuck|smut|insert(?:ed)? himself|\bhis member|\bgrope|\bthrust|\bpant(?:y|ies)\b|\bloli\b|\bescort\b|bondage|fisting|nipple|yaoi|moaning|\btwink\b|(?:mother)?fuck(?:er)?|ass ?hole|\bass\b|\bbutt\b|\bcondom\b|turned(?:-| )on|\blick(?:s|ed|ing)?|\bsuck(?:s|ed|ing)?|\bstroke?(?:s|d|ing)?|\bfuk\b|tease(?:d)? (?:her|him)|\bperv|(?:dirty talk|talk dirty)|\bher (?:throat|mouth)\b|\bunzip|\bunbutton|\bpink nub\b|\bclimax|\bf(?:v|\*)ck|shit|bitch|skank|dumbass|virgin|marijuana|cocaine|xhamster|beeg.com|xvideos|redtube|xnxx|tube8|spankwire|fapdu|empflix|brazzers`)
}

// Censor removes censored words in a text
func Censor(text *string) {
	if text == nil {
		return
	}
	re := GetCensorRegex()
	*text = re.ReplaceAllString(*text, "")
}

// BadWordFound determines if a censored
// word is found in a text.
func BadWordFound(text *string) bool {
	if text == nil {
		return false
	}
	re := GetCensorRegex()
	substr := re.FindString(*text)
	if substr != "" {
		return true
	}
	return false
}
