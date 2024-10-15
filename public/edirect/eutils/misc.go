// ===========================================================================
//
//                            PUBLIC DOMAIN NOTICE
//            National Center for Biotechnology Information (NCBI)
//
//  This software/database is a "United States Government Work" under the
//  terms of the United States Copyright Act. It was written as part of
//  the author's official duties as a United States Government employee and
//  thus cannot be copyrighted. This software/database is freely available
//  to the public for use. The National Library of Medicine and the U.S.
//  Government do not place any restriction on its use or reproduction.
//  We would, however, appreciate having the NCBI and the author cited in
//  any work or product based on this material.
//
//  Although all reasonable efforts have been taken to ensure the accuracy
//  and reliability of the software and data, the NLM and the U.S.
//  Government do not and cannot warrant the performance or results that
//  may be obtained by using this software or data. The NLM and the U.S.
//  Government disclaim all warranties, express or implied, including
//  warranties of performance, merchantability or fitness for any particular
//  purpose.
//
// ===========================================================================
//
// File Name:  misc.go
//
// Author:  Jonathan Kans
//
// ==========================================================================

package eutils

import (
	"cmp"
	"encoding/hex"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"html"
	"slices"
	"strconv"
	"strings"
	"sync"
	"unicode"
)

// constants for markup policy
const (
	NOMARKUP = iota
	FUSE
	SPACE
	PERIOD
	CONCISE
	BRACKETS
	MARKDOWN
	SLASH
	TAGS
	TERSE
)

// removes mixed content tags
var mfix *strings.Replacer

// lower-case version of combiningAccents
var lcBadAccents map[string]bool

var invisibleRunes = map[rune]bool{
	0x00A0: true, // NO-BREAK SPACE
	0x00AD: true, // SOFT HYPHEN
	0x034F: true, // COMBINING GRAPHEME JOINER
	0x061C: true, // ARABIC LETTER MARK
	0x115F: true, // HANGUL CHOSEONG FILLER
	0x1160: true, // HANGUL JUNGSEONG FILLER
	0x17B4: true, // KHMER VOWEL INHERENT AQ
	0x17B5: true, // KHMER VOWEL INHERENT AA
	0x180E: true, // MONGOLIAN VOWEL SEPARATOR
	0x2000: true, // EN QUAD
	0x2001: true, // EM QUAD
	0x2002: true, // EN SPACE
	0x2003: true, // EM SPACE
	0x2004: true, // THREE-PER-EM SPACE
	0x2005: true, // FOUR-PER-EM SPACE
	0x2006: true, // SIX-PER-EM SPACE
	0x2007: true, // FIGURE SPACE
	0x2008: true, // PUNCTUATION SPACE
	0x2009: true, // THIN SPACE
	0x200A: true, // HAIR SPACE
	0x200B: true, // ZERO WIDTH SPACE
	0x200C: true, // ZERO WIDTH NON-JOINER
	0x200D: true, // ZERO WIDTH JOINER
	0x200E: true, // LEFT-TO-RIGHT MARK
	0x200F: true, // RIGHT-TO-LEFT MARK
	0x202F: true, // NARROW NO-BREAK SPACE
	0x205F: true, // MEDIUM MATHEMATICAL SPACE
	0x2060: true, // WORD JOINER
	0x2061: true, // FUNCTION APPLICATION
	0x2062: true, // INVISIBLE TIMES
	0x2063: true, // INVISIBLE SEPARATOR
	0x2064: true, // INVISIBLE PLUS
	0x206A: true, // INHIBIT SYMMETRIC SWAPPING
	0x206B: true, // ACTIVATE SYMMETRIC SWAPPING
	0x206C: true, // INHIBIT ARABIC FORM SHAPING
	0x206D: true, // ACTIVATE ARABIC FORM SHAPING
	0x206E: true, // NATIONAL DIGIT SHAPES
	0x206F: true, // NOMINAL DIGIT SHAPES
	0x3000: true, // IDEOGRAPHIC SPACE
	0x2800: true, // BRAILLE PATTERN BLANK
	0x3164: true, // HANGUL FILLER
	0xFEFF: true, // ZERO WIDTH NO-BREAK SPACE
	0xFFA0: true, // HALFWIDTH HANGUL FILLER
}

var subScriptRunes = map[rune]rune{
	0x2080: '0',
	0x2081: '1',
	0x2082: '2',
	0x2083: '3',
	0x2084: '4',
	0x2085: '5',
	0x2086: '6',
	0x2087: '7',
	0x2088: '8',
	0x2089: '9',
	0x208A: '+',
	0x208B: '-',
	0x208C: '=',
	0x208D: '(',
	0x208E: ')',
	0x2090: 'a',
	0x2091: 'e',
	0x2092: 'o',
	0x2093: 'x',
	0x2094: 'e',
	0x2095: 'h',
	0x2096: 'k',
	0x2097: 'l',
	0x2098: 'm',
	0x2099: 'n',
	0x209A: 'p',
	0x209B: 's',
	0x209C: 't',
}

var superScriptRunes = map[rune]rune{
	0x00B2: '2',
	0x00B3: '3',
	0x00B9: '1',
	0x2070: '0',
	0x2071: '1',
	0x2074: '4',
	0x2075: '5',
	0x2076: '6',
	0x2077: '7',
	0x2078: '8',
	0x2079: '9',
	0x207A: '+',
	0x207B: '-',
	0x207C: '=',
	0x207D: '(',
	0x207E: ')',
	0x207F: 'n',
}

var isStopWord = map[string]bool{
	"a":             true,
	"about":         true,
	"above":         true,
	"abs":           true,
	"accordingly":   true,
	"across":        true,
	"after":         true,
	"afterwards":    true,
	"again":         true,
	"against":       true,
	"all":           true,
	"almost":        true,
	"alone":         true,
	"along":         true,
	"already":       true,
	"also":          true,
	"although":      true,
	"always":        true,
	"am":            true,
	"among":         true,
	"amongst":       true,
	"an":            true,
	"analyze":       true,
	"and":           true,
	"another":       true,
	"any":           true,
	"anyhow":        true,
	"anyone":        true,
	"anything":      true,
	"anywhere":      true,
	"applicable":    true,
	"apply":         true,
	"are":           true,
	"arise":         true,
	"around":        true,
	"as":            true,
	"assume":        true,
	"at":            true,
	"be":            true,
	"became":        true,
	"because":       true,
	"become":        true,
	"becomes":       true,
	"becoming":      true,
	"been":          true,
	"before":        true,
	"beforehand":    true,
	"being":         true,
	"below":         true,
	"beside":        true,
	"besides":       true,
	"between":       true,
	"beyond":        true,
	"both":          true,
	"but":           true,
	"by":            true,
	"came":          true,
	"can":           true,
	"cannot":        true,
	"cc":            true,
	"cm":            true,
	"come":          true,
	"compare":       true,
	"could":         true,
	"de":            true,
	"dealing":       true,
	"department":    true,
	"depend":        true,
	"did":           true,
	"discover":      true,
	"dl":            true,
	"do":            true,
	"does":          true,
	"done":          true,
	"due":           true,
	"during":        true,
	"each":          true,
	"ec":            true,
	"ed":            true,
	"effected":      true,
	"eg":            true,
	"either":        true,
	"else":          true,
	"elsewhere":     true,
	"enough":        true,
	"especially":    true,
	"et":            true,
	"etc":           true,
	"ever":          true,
	"every":         true,
	"everyone":      true,
	"everything":    true,
	"everywhere":    true,
	"except":        true,
	"find":          true,
	"for":           true,
	"found":         true,
	"from":          true,
	"further":       true,
	"gave":          true,
	"get":           true,
	"give":          true,
	"go":            true,
	"gone":          true,
	"got":           true,
	"gov":           true,
	"had":           true,
	"has":           true,
	"have":          true,
	"having":        true,
	"he":            true,
	"hence":         true,
	"her":           true,
	"here":          true,
	"hereafter":     true,
	"hereby":        true,
	"herein":        true,
	"hereupon":      true,
	"hers":          true,
	"herself":       true,
	"him":           true,
	"himself":       true,
	"his":           true,
	"how":           true,
	"however":       true,
	"hr":            true,
	"i":             true,
	"ie":            true,
	"if":            true,
	"ii":            true,
	"iii":           true,
	"immediately":   true,
	"importance":    true,
	"important":     true,
	"in":            true,
	"inc":           true,
	"incl":          true,
	"indeed":        true,
	"into":          true,
	"investigate":   true,
	"is":            true,
	"it":            true,
	"its":           true,
	"itself":        true,
	"just":          true,
	"keep":          true,
	"kept":          true,
	"kg":            true,
	"km":            true,
	"last":          true,
	"latter":        true,
	"latterly":      true,
	"lb":            true,
	"ld":            true,
	"letter":        true,
	"like":          true,
	"ltd":           true,
	"made":          true,
	"mainly":        true,
	"make":          true,
	"many":          true,
	"may":           true,
	"me":            true,
	"meanwhile":     true,
	"mg":            true,
	"might":         true,
	"ml":            true,
	"mm":            true,
	"mo":            true,
	"more":          true,
	"moreover":      true,
	"most":          true,
	"mostly":        true,
	"mr":            true,
	"much":          true,
	"mug":           true,
	"must":          true,
	"my":            true,
	"myself":        true,
	"namely":        true,
	"nearly":        true,
	"necessarily":   true,
	"neither":       true,
	"never":         true,
	"nevertheless":  true,
	"next":          true,
	"no":            true,
	"nobody":        true,
	"noone":         true,
	"nor":           true,
	"normally":      true,
	"nos":           true,
	"not":           true,
	"noted":         true,
	"nothing":       true,
	"now":           true,
	"nowhere":       true,
	"obtained":      true,
	"of":            true,
	"off":           true,
	"often":         true,
	"on":            true,
	"only":          true,
	"onto":          true,
	"or":            true,
	"other":         true,
	"others":        true,
	"otherwise":     true,
	"ought":         true,
	"our":           true,
	"ours":          true,
	"ourselves":     true,
	"out":           true,
	"over":          true,
	"overall":       true,
	"owing":         true,
	"own":           true,
	"oz":            true,
	"particularly":  true,
	"per":           true,
	"perhaps":       true,
	"pm":            true,
	"pmid":          true,
	"precede":       true,
	"predominantly": true,
	"present":       true,
	"presently":     true,
	"previously":    true,
	"primarily":     true,
	"promptly":      true,
	"pt":            true,
	"quickly":       true,
	"quite":         true,
	"quot":          true,
	"rather":        true,
	"readily":       true,
	"really":        true,
	"recently":      true,
	"refs":          true,
	"regarding":     true,
	"relate":        true,
	"said":          true,
	"same":          true,
	"seem":          true,
	"seemed":        true,
	"seeming":       true,
	"seems":         true,
	"seen":          true,
	"seriously":     true,
	"several":       true,
	"shall":         true,
	"she":           true,
	"should":        true,
	"show":          true,
	"showed":        true,
	"shown":         true,
	"shows":         true,
	"significantly": true,
	"since":         true,
	"slightly":      true,
	"so":            true,
	"some":          true,
	"somehow":       true,
	"someone":       true,
	"something":     true,
	"sometime":      true,
	"sometimes":     true,
	"somewhat":      true,
	"somewhere":     true,
	"soon":          true,
	"specifically":  true,
	"still":         true,
	"strongly":      true,
	"studied":       true,
	"studies":       true,
	"study":         true,
	"sub":           true,
	"substantially": true,
	"such":          true,
	"sufficiently":  true,
	"take":          true,
	"tell":          true,
	"th":            true,
	"than":          true,
	"that":          true,
	"the":           true,
	"their":         true,
	"theirs":        true,
	"them":          true,
	"themselves":    true,
	"then":          true,
	"thence":        true,
	"there":         true,
	"thereafter":    true,
	"thereby":       true,
	"therefore":     true,
	"therein":       true,
	"thereupon":     true,
	"these":         true,
	"they":          true,
	"this":          true,
	"thorough":      true,
	"those":         true,
	"though":        true,
	"through":       true,
	"throughout":    true,
	"thru":          true,
	"thus":          true,
	"to":            true,
	"together":      true,
	"too":           true,
	"toward":        true,
	"towards":       true,
	"try":           true,
	"type":          true,
	"ug":            true,
	"under":         true,
	"unless":        true,
	"until":         true,
	"up":            true,
	"upon":          true,
	"us":            true,
	"use":           true,
	"used":          true,
	"usefully":      true,
	"usefulness":    true,
	"using":         true,
	"usually":       true,
	"various":       true,
	"very":          true,
	"via":           true,
	"was":           true,
	"we":            true,
	"were":          true,
	"what":          true,
	"whatever":      true,
	"when":          true,
	"whence":        true,
	"whenever":      true,
	"where":         true,
	"whereafter":    true,
	"whereas":       true,
	"whereby":       true,
	"wherein":       true,
	"whereupon":     true,
	"wherever":      true,
	"whether":       true,
	"which":         true,
	"while":         true,
	"whither":       true,
	"who":           true,
	"whoever":       true,
	"whom":          true,
	"whose":         true,
	"why":           true,
	"will":          true,
	"with":          true,
	"within":        true,
	"without":       true,
	"wk":            true,
	"would":         true,
	"wt":            true,
	"yet":           true,
	"you":           true,
	"your":          true,
	"yours":         true,
	"yourself":      true,
	"yourselves":    true,
	"yr":            true,
}

var htmlRepair = map[string]string{
	"&amp;lt;b&amp;gt;":     "<b>",
	"&amp;lt;i&amp;gt;":     "<i>",
	"&amp;lt;u&amp;gt;":     "<u>",
	"&amp;lt;/b&amp;gt;":    "</b>",
	"&amp;lt;/i&amp;gt;":    "</i>",
	"&amp;lt;/u&amp;gt;":    "</u>",
	"&amp;lt;b/&amp;gt;":    "<b/>",
	"&amp;lt;i/&amp;gt;":    "<i/>",
	"&amp;lt;u/&amp;gt;":    "<u/>",
	"&amp;lt;b /&amp;gt;":   "<b/>",
	"&amp;lt;i /&amp;gt;":   "<i/>",
	"&amp;lt;u /&amp;gt;":   "<u/>",
	"&amp;lt;sub&amp;gt;":   "<sub>",
	"&amp;lt;sup&amp;gt;":   "<sup>",
	"&amp;lt;/sub&amp;gt;":  "</sub>",
	"&amp;lt;/sup&amp;gt;":  "</sup>",
	"&amp;lt;sub/&amp;gt;":  "<sub/>",
	"&amp;lt;sup/&amp;gt;":  "<sup/>",
	"&amp;lt;sub /&amp;gt;": "<sub/>",
	"&amp;lt;sup /&amp;gt;": "<sup/>",
	"&lt;b&gt;":             "<b>",
	"&lt;i&gt;":             "<i>",
	"&lt;u&gt;":             "<u>",
	"&lt;/b&gt;":            "</b>",
	"&lt;/i&gt;":            "</i>",
	"&lt;/u&gt;":            "</u>",
	"&lt;b/&gt;":            "<b/>",
	"&lt;i/&gt;":            "<i/>",
	"&lt;u/&gt;":            "<u/>",
	"&lt;b /&gt;":           "<b/>",
	"&lt;i /&gt;":           "<i/>",
	"&lt;u /&gt;":           "<u/>",
	"&lt;sub&gt;":           "<sub>",
	"&lt;sup&gt;":           "<sup>",
	"&lt;/sub&gt;":          "</sub>",
	"&lt;/sup&gt;":          "</sup>",
	"&lt;sub/&gt;":          "<sub/>",
	"&lt;sup/&gt;":          "<sup/>",
	"&lt;sub /&gt;":         "<sub/>",
	"&lt;sup /&gt;":         "<sup/>",
}

var hyphenatedPrefixes = map[string]bool{
	"anti":    true,
	"bi":      true,
	"co":      true,
	"contra":  true,
	"counter": true,
	"de":      true,
	"di":      true,
	"extra":   true,
	"infra":   true,
	"inter":   true,
	"intra":   true,
	"micro":   true,
	"mid":     true,
	"mono":    true,
	"multi":   true,
	"non":     true,
	"over":    true,
	"peri":    true,
	"post":    true,
	"pre":     true,
	"pro":     true,
	"proto":   true,
	"pseudo":  true,
	"re":      true,
	"semi":    true,
	"sub":     true,
	"super":   true,
	"supra":   true,
	"tetra":   true,
	"trans":   true,
	"tri":     true,
	"ultra":   true,
	"un":      true,
	"under":   true,
	"whole":   true,
}

var primedPrefixes = map[string]bool{
	"5": true,
	"3": true,
}

var primedSuffix = map[string]bool{
	"s": true,
}

var revComp = map[rune]rune{
	'A': 'T',
	'B': 'V',
	'C': 'G',
	'D': 'H',
	'G': 'C',
	'H': 'D',
	'K': 'M',
	'M': 'K',
	'N': 'N',
	'R': 'Y',
	'S': 'S',
	'T': 'A',
	'U': 'A',
	'V': 'B',
	'W': 'W',
	'X': 'X',
	'Y': 'R',
	'a': 't',
	'b': 'v',
	'c': 'g',
	'd': 'h',
	'g': 'c',
	'h': 'd',
	'k': 'm',
	'm': 'k',
	'n': 'n',
	'r': 'y',
	's': 's',
	't': 'a',
	'u': 'a',
	'v': 'b',
	'w': 'w',
	'x': 'x',
	'y': 'r',
}

var aaTo3 = map[string]string{
	"*": "Ter",
	"-": "Gap",
	"A": "Ala",
	"B": "Asx",
	"C": "Cys",
	"D": "Asp",
	"E": "Glu",
	"F": "Phe",
	"G": "Gly",
	"H": "His",
	"I": "Ile",
	"J": "Xle",
	"K": "Lys",
	"L": "Leu",
	"M": "Met",
	"N": "Asn",
	"O": "Pyl",
	"P": "Pro",
	"Q": "Gln",
	"R": "Arg",
	"S": "Ser",
	"T": "Thr",
	"U": "Sec",
	"V": "Val",
	"W": "Trp",
	"X": "Xxx",
	"Y": "Tyr",
	"Z": "Glx",
	"a": "Ala",
	"b": "Asx",
	"c": "Cys",
	"d": "Asp",
	"e": "Glu",
	"f": "Phe",
	"g": "Gly",
	"h": "His",
	"i": "Ile",
	"j": "Xle",
	"k": "Lys",
	"l": "Leu",
	"m": "Met",
	"n": "Asn",
	"o": "Pyl",
	"p": "Pro",
	"q": "Gln",
	"r": "Arg",
	"s": "Ser",
	"t": "Thr",
	"u": "Sec",
	"v": "Val",
	"w": "Trp",
	"x": "Xxx",
	"y": "Tyr",
	"z": "Glx",
}

var aaTo1 = map[string]string{
	"*":   "*",
	"-":   "-",
	"Ala": "A",
	"Arg": "R",
	"Asn": "N",
	"Asp": "D",
	"Asx": "B",
	"Cys": "C",
	"Gap": "-",
	"Gln": "Q",
	"Glu": "E",
	"Glx": "Z",
	"Gly": "G",
	"His": "H",
	"Ile": "I",
	"Leu": "L",
	"Lys": "K",
	"Met": "M",
	"Phe": "F",
	"Pro": "P",
	"Pyl": "O",
	"Sec": "U",
	"Ser": "S",
	"Ter": "*",
	"Thr": "T",
	"Trp": "W",
	"Tyr": "Y",
	"Val": "V",
	"Xle": "J",
	"Xxx": "X",
	"ala": "A",
	"arg": "R",
	"asn": "N",
	"asp": "D",
	"asx": "B",
	"cys": "C",
	"gap": "-",
	"gln": "Q",
	"glu": "E",
	"glx": "Z",
	"gly": "G",
	"his": "H",
	"ile": "I",
	"leu": "L",
	"lys": "K",
	"met": "M",
	"phe": "F",
	"pro": "P",
	"pyl": "O",
	"sec": "U",
	"ser": "S",
	"ter": "*",
	"thr": "T",
	"trp": "W",
	"tyr": "Y",
	"val": "V",
	"xle": "J",
	"xxx": "X",
}

var expandNuc = map[string]string{
	"A": "A",
	"B": "CGT",
	"C": "C",
	"D": "AGT",
	"G": "G",
	"H": "ACT",
	"K": "GT",
	"M": "AC",
	"N": "ACGT",
	"R": "AG",
	"S": "GC",
	"T": "T",
	"V": "ACG",
	"W": "AT",
	"X": "ACGT",
	"Y": "CT",
	"a": "a",
	"b": "cgt",
	"c": "c",
	"d": "agt",
	"g": "g",
	"h": "act",
	"k": "gt",
	"m": "ac",
	"n": "acgt",
	"r": "ag",
	"s": "gc",
	"t": "t",
	"v": "acg",
	"w": "at",
	"x": "acgt",
	"y": "ct",
}

/*
	var conv = []string{"A", "C", "G", "T"}
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			for k := 0; k < 4; k++ {
				for l := 0; l < 4; l++ {
					base := conv[i] + conv[j] + conv[k] + conv[l]
					idx := i*64 + j*16 + k*4 + l
					fmt.Fprintf(os.Stdout, "\t%d: \"%s\",\n", idx, base)
				}
			}
		}
	}
*/

var ncbi2naToIupac = map[int]string{
	0:   "AAAA",
	1:   "AAAC",
	2:   "AAAG",
	3:   "AAAT",
	4:   "AACA",
	5:   "AACC",
	6:   "AACG",
	7:   "AACT",
	8:   "AAGA",
	9:   "AAGC",
	10:  "AAGG",
	11:  "AAGT",
	12:  "AATA",
	13:  "AATC",
	14:  "AATG",
	15:  "AATT",
	16:  "ACAA",
	17:  "ACAC",
	18:  "ACAG",
	19:  "ACAT",
	20:  "ACCA",
	21:  "ACCC",
	22:  "ACCG",
	23:  "ACCT",
	24:  "ACGA",
	25:  "ACGC",
	26:  "ACGG",
	27:  "ACGT",
	28:  "ACTA",
	29:  "ACTC",
	30:  "ACTG",
	31:  "ACTT",
	32:  "AGAA",
	33:  "AGAC",
	34:  "AGAG",
	35:  "AGAT",
	36:  "AGCA",
	37:  "AGCC",
	38:  "AGCG",
	39:  "AGCT",
	40:  "AGGA",
	41:  "AGGC",
	42:  "AGGG",
	43:  "AGGT",
	44:  "AGTA",
	45:  "AGTC",
	46:  "AGTG",
	47:  "AGTT",
	48:  "ATAA",
	49:  "ATAC",
	50:  "ATAG",
	51:  "ATAT",
	52:  "ATCA",
	53:  "ATCC",
	54:  "ATCG",
	55:  "ATCT",
	56:  "ATGA",
	57:  "ATGC",
	58:  "ATGG",
	59:  "ATGT",
	60:  "ATTA",
	61:  "ATTC",
	62:  "ATTG",
	63:  "ATTT",
	64:  "CAAA",
	65:  "CAAC",
	66:  "CAAG",
	67:  "CAAT",
	68:  "CACA",
	69:  "CACC",
	70:  "CACG",
	71:  "CACT",
	72:  "CAGA",
	73:  "CAGC",
	74:  "CAGG",
	75:  "CAGT",
	76:  "CATA",
	77:  "CATC",
	78:  "CATG",
	79:  "CATT",
	80:  "CCAA",
	81:  "CCAC",
	82:  "CCAG",
	83:  "CCAT",
	84:  "CCCA",
	85:  "CCCC",
	86:  "CCCG",
	87:  "CCCT",
	88:  "CCGA",
	89:  "CCGC",
	90:  "CCGG",
	91:  "CCGT",
	92:  "CCTA",
	93:  "CCTC",
	94:  "CCTG",
	95:  "CCTT",
	96:  "CGAA",
	97:  "CGAC",
	98:  "CGAG",
	99:  "CGAT",
	100: "CGCA",
	101: "CGCC",
	102: "CGCG",
	103: "CGCT",
	104: "CGGA",
	105: "CGGC",
	106: "CGGG",
	107: "CGGT",
	108: "CGTA",
	109: "CGTC",
	110: "CGTG",
	111: "CGTT",
	112: "CTAA",
	113: "CTAC",
	114: "CTAG",
	115: "CTAT",
	116: "CTCA",
	117: "CTCC",
	118: "CTCG",
	119: "CTCT",
	120: "CTGA",
	121: "CTGC",
	122: "CTGG",
	123: "CTGT",
	124: "CTTA",
	125: "CTTC",
	126: "CTTG",
	127: "CTTT",
	128: "GAAA",
	129: "GAAC",
	130: "GAAG",
	131: "GAAT",
	132: "GACA",
	133: "GACC",
	134: "GACG",
	135: "GACT",
	136: "GAGA",
	137: "GAGC",
	138: "GAGG",
	139: "GAGT",
	140: "GATA",
	141: "GATC",
	142: "GATG",
	143: "GATT",
	144: "GCAA",
	145: "GCAC",
	146: "GCAG",
	147: "GCAT",
	148: "GCCA",
	149: "GCCC",
	150: "GCCG",
	151: "GCCT",
	152: "GCGA",
	153: "GCGC",
	154: "GCGG",
	155: "GCGT",
	156: "GCTA",
	157: "GCTC",
	158: "GCTG",
	159: "GCTT",
	160: "GGAA",
	161: "GGAC",
	162: "GGAG",
	163: "GGAT",
	164: "GGCA",
	165: "GGCC",
	166: "GGCG",
	167: "GGCT",
	168: "GGGA",
	169: "GGGC",
	170: "GGGG",
	171: "GGGT",
	172: "GGTA",
	173: "GGTC",
	174: "GGTG",
	175: "GGTT",
	176: "GTAA",
	177: "GTAC",
	178: "GTAG",
	179: "GTAT",
	180: "GTCA",
	181: "GTCC",
	182: "GTCG",
	183: "GTCT",
	184: "GTGA",
	185: "GTGC",
	186: "GTGG",
	187: "GTGT",
	188: "GTTA",
	189: "GTTC",
	190: "GTTG",
	191: "GTTT",
	192: "TAAA",
	193: "TAAC",
	194: "TAAG",
	195: "TAAT",
	196: "TACA",
	197: "TACC",
	198: "TACG",
	199: "TACT",
	200: "TAGA",
	201: "TAGC",
	202: "TAGG",
	203: "TAGT",
	204: "TATA",
	205: "TATC",
	206: "TATG",
	207: "TATT",
	208: "TCAA",
	209: "TCAC",
	210: "TCAG",
	211: "TCAT",
	212: "TCCA",
	213: "TCCC",
	214: "TCCG",
	215: "TCCT",
	216: "TCGA",
	217: "TCGC",
	218: "TCGG",
	219: "TCGT",
	220: "TCTA",
	221: "TCTC",
	222: "TCTG",
	223: "TCTT",
	224: "TGAA",
	225: "TGAC",
	226: "TGAG",
	227: "TGAT",
	228: "TGCA",
	229: "TGCC",
	230: "TGCG",
	231: "TGCT",
	232: "TGGA",
	233: "TGGC",
	234: "TGGG",
	235: "TGGT",
	236: "TGTA",
	237: "TGTC",
	238: "TGTG",
	239: "TGTT",
	240: "TTAA",
	241: "TTAC",
	242: "TTAG",
	243: "TTAT",
	244: "TTCA",
	245: "TTCC",
	246: "TTCG",
	247: "TTCT",
	248: "TTGA",
	249: "TTGC",
	250: "TTGG",
	251: "TTGT",
	252: "TTTA",
	253: "TTTC",
	254: "TTTG",
	255: "TTTT",
}

/*
	var conv = []string{"N", "A", "C", "M", "G", "R", "S", "V", "T", "W", "Y", "H", "K", "D", "B", "N"}
	for i := 0; i < 16; i++ {
		for j := 0; j < 16; j++ {
			base := conv[i] + conv[j]
			idx := i*16 + j
			fmt.Fprintf(os.Stdout, "\t%d: \"%s\",\n", idx, base)
		}
	}
*/

var ncbi4naToIupac = map[int]string{
	0:   "NN",
	1:   "NA",
	2:   "NC",
	3:   "NM",
	4:   "NG",
	5:   "NR",
	6:   "NS",
	7:   "NV",
	8:   "NT",
	9:   "NW",
	10:  "NY",
	11:  "NH",
	12:  "NK",
	13:  "ND",
	14:  "NB",
	15:  "NN",
	16:  "AN",
	17:  "AA",
	18:  "AC",
	19:  "AM",
	20:  "AG",
	21:  "AR",
	22:  "AS",
	23:  "AV",
	24:  "AT",
	25:  "AW",
	26:  "AY",
	27:  "AH",
	28:  "AK",
	29:  "AD",
	30:  "AB",
	31:  "AN",
	32:  "CN",
	33:  "CA",
	34:  "CC",
	35:  "CM",
	36:  "CG",
	37:  "CR",
	38:  "CS",
	39:  "CV",
	40:  "CT",
	41:  "CW",
	42:  "CY",
	43:  "CH",
	44:  "CK",
	45:  "CD",
	46:  "CB",
	47:  "CN",
	48:  "MN",
	49:  "MA",
	50:  "MC",
	51:  "MM",
	52:  "MG",
	53:  "MR",
	54:  "MS",
	55:  "MV",
	56:  "MT",
	57:  "MW",
	58:  "MY",
	59:  "MH",
	60:  "MK",
	61:  "MD",
	62:  "MB",
	63:  "MN",
	64:  "GN",
	65:  "GA",
	66:  "GC",
	67:  "GM",
	68:  "GG",
	69:  "GR",
	70:  "GS",
	71:  "GV",
	72:  "GT",
	73:  "GW",
	74:  "GY",
	75:  "GH",
	76:  "GK",
	77:  "GD",
	78:  "GB",
	79:  "GN",
	80:  "RN",
	81:  "RA",
	82:  "RC",
	83:  "RM",
	84:  "RG",
	85:  "RR",
	86:  "RS",
	87:  "RV",
	88:  "RT",
	89:  "RW",
	90:  "RY",
	91:  "RH",
	92:  "RK",
	93:  "RD",
	94:  "RB",
	95:  "RN",
	96:  "SN",
	97:  "SA",
	98:  "SC",
	99:  "SM",
	100: "SG",
	101: "SR",
	102: "SS",
	103: "SV",
	104: "ST",
	105: "SW",
	106: "SY",
	107: "SH",
	108: "SK",
	109: "SD",
	110: "SB",
	111: "SN",
	112: "VN",
	113: "VA",
	114: "VC",
	115: "VM",
	116: "VG",
	117: "VR",
	118: "VS",
	119: "VV",
	120: "VT",
	121: "VW",
	122: "VY",
	123: "VH",
	124: "VK",
	125: "VD",
	126: "VB",
	127: "VN",
	128: "TN",
	129: "TA",
	130: "TC",
	131: "TM",
	132: "TG",
	133: "TR",
	134: "TS",
	135: "TV",
	136: "TT",
	137: "TW",
	138: "TY",
	139: "TH",
	140: "TK",
	141: "TD",
	142: "TB",
	143: "TN",
	144: "WN",
	145: "WA",
	146: "WC",
	147: "WM",
	148: "WG",
	149: "WR",
	150: "WS",
	151: "WV",
	152: "WT",
	153: "WW",
	154: "WY",
	155: "WH",
	156: "WK",
	157: "WD",
	158: "WB",
	159: "WN",
	160: "YN",
	161: "YA",
	162: "YC",
	163: "YM",
	164: "YG",
	165: "YR",
	166: "YS",
	167: "YV",
	168: "YT",
	169: "YW",
	170: "YY",
	171: "YH",
	172: "YK",
	173: "YD",
	174: "YB",
	175: "YN",
	176: "HN",
	177: "HA",
	178: "HC",
	179: "HM",
	180: "HG",
	181: "HR",
	182: "HS",
	183: "HV",
	184: "HT",
	185: "HW",
	186: "HY",
	187: "HH",
	188: "HK",
	189: "HD",
	190: "HB",
	191: "HN",
	192: "KN",
	193: "KA",
	194: "KC",
	195: "KM",
	196: "KG",
	197: "KR",
	198: "KS",
	199: "KV",
	200: "KT",
	201: "KW",
	202: "KY",
	203: "KH",
	204: "KK",
	205: "KD",
	206: "KB",
	207: "KN",
	208: "DN",
	209: "DA",
	210: "DC",
	211: "DM",
	212: "DG",
	213: "DR",
	214: "DS",
	215: "DV",
	216: "DT",
	217: "DW",
	218: "DY",
	219: "DH",
	220: "DK",
	221: "DD",
	222: "DB",
	223: "DN",
	224: "BN",
	225: "BA",
	226: "BC",
	227: "BM",
	228: "BG",
	229: "BR",
	230: "BS",
	231: "BV",
	232: "BT",
	233: "BW",
	234: "BY",
	235: "BH",
	236: "BK",
	237: "BD",
	238: "BB",
	239: "BN",
	240: "NN",
	241: "NA",
	242: "NC",
	243: "NM",
	244: "NG",
	245: "NR",
	246: "NS",
	247: "NV",
	248: "NT",
	249: "NW",
	250: "NY",
	251: "NH",
	252: "NK",
	253: "ND",
	254: "NB",
	255: "NN",
}

// use RelaxString to convert non-alphanumeric characters to spaces

// CleanupAuthor fixes misused letters and accents
func CleanupAuthor(str string, wrp bool) string {

	if str == "" {
		return str
	}

	str = strings.Replace(str, "\n", " ", -1)

	str = FixMisusedLetters(str, false, true, false)
	str = TransformAccents(str, false, false)
	// convert numeric encoding to apostrophe
	str = strings.Replace(str, "&#39;", "'", -1)
	// remove space following apostrophe
	str = strings.Replace(str, "' ", "'", -1)
	// but leave apostrophe present in result

	if HasUnicodeMarkup(str) {
		str = RepairUnicodeMarkup(str, SPACE)
	}
	if HasAngleOrAmpersandEncoding(str) {
		str = RepairTableMarkup(str, SPACE)
		str = RemoveHTMLDecorations(str)
		if wrp {
			str = encodeAngleBracketsAndAmpersand(str)
		}
	}
	if HasBadSpace(str) {
		str = CleanupBadSpaces(str)
	}
	if HasAdjacentSpaces(str) {
		str = CompressRunsOfSpaces(str)
	}
	if HasExtraSpaces(str) {
		str = RemoveExtraSpaces(str)
	}

	return str
}

// CleanupBadSpaces converts non-ASCII Unicode spaces to ASCII space
func CleanupBadSpaces(str string) string {

	var buffer strings.Builder

	for _, ch := range str {
		if ch < 128 {
			buffer.WriteRune(ch)
			continue
		}

		if unicode.IsSpace(ch) {
			buffer.WriteRune(' ')
		} else if ch >= 0x0080 && ch <= 0x009F {
			// absorb control characters
		} else if ch == 0x00A0 {
			// no-break space
			buffer.WriteRune(' ')
		} else if ch >= 0x2002 && ch <= 0x200B {
			// special spaces (en space through zero width space)
			buffer.WriteRune(' ')
		} else if ch == 0x202F {
			// narrow no-break space
			buffer.WriteRune(' ')
		} else if ch == 0x3000 {
			// ideographic space
			buffer.WriteRune(' ')
		} else if ch >= 0xE000 && ch <= 0xF8FF {
			// unassigned values (private use)
			buffer.WriteRune('(')
			buffer.WriteRune('?')
			buffer.WriteRune(')')
		} else if ch == 0xFEFF {
			// zero width no-break space
			buffer.WriteRune(' ')
		} else {
			buffer.WriteRune(ch)
		}
	}

	return buffer.String()
}

// CleanupBasic convert superscripts and subscripts
func CleanupBasic(str string, wrp bool) string {

	if str == "" {
		return str
	}

	str = strings.Replace(str, "\n", " ", -1)

	if HasUnicodeMarkup(str) {
		str = RepairUnicodeMarkup(str, SPACE)
	}
	if HasAngleOrAmpersandEncoding(str) {
		str = RepairTableMarkup(str, SPACE)
		str = RemoveHTMLDecorations(str)
		if wrp {
			str = encodeAngleBracketsAndAmpersand(str)
		}
	}
	if HasBadSpace(str) {
		str = CleanupBadSpaces(str)
	}
	if HasAdjacentSpaces(str) {
		str = CompressRunsOfSpaces(str)
	}
	if HasExtraSpaces(str) {
		str = RemoveExtraSpaces(str)
	}

	return str
}

// CleanupContents performs optional operations on XML content strings
func CleanupContents(str string, ascii, amper, mixed bool) string {

	if doCompress {
		if !allowEmbed {
			if ascii && HasBadSpace(str) {
				str = CleanupBadSpaces(str)
			}
		}
		if HasAdjacentSpacesOrNewline(str) {
			str = CompressRunsOfSpaces(str)
		}
	}
	if allowEmbed {
		if amper {
			str = RepairEncodedMarkup(str)
		}
	}
	if doScript {
		if mixed && HasAngleOrAmpersandEncoding(str) {
			str = RepairScriptMarkup(str, CONCISE)
		}
	}
	if doMathML {
		if mixed && HasAngleOrAmpersandEncoding(str) {
			str = RepairMathMLMarkup(str, mathMLFix)
		}
	}
	if doStrict {
		if ascii {
			if HasUnicodeMarkup(str) {
				str = RepairUnicodeMarkup(str, TAGS)
				mixed = true
			}
		}
		if mixed || amper {
			if HasAngleOrAmpersandEncoding(str) {
				str = RepairScriptMarkup(str, CONCISE)
				str = RepairTableMarkup(str, SPACE)
				// call RepairScriptMarkup before RemoveEmbeddedMarkup
				str = RemoveEmbeddedMarkup(str)
			}
		}
		if ascii && HasBadSpace(str) {
			str = CleanupBadSpaces(str)
		}
		if HasAdjacentSpaces(str) {
			str = CompressRunsOfSpaces(str)
		}
		// Remove MathML artifact
		if NeedsTightening(str) {
			str = TightenParentheses(str)
		}
	}
	if doMixed {
		if ascii {
			if HasUnicodeMarkup(str) {
				str = RepairUnicodeMarkup(str, TAGS)
				mixed = true
			}
		}
		if mixed {
			str = DoTrimFlankingHTML(str)
		}
		if ascii && HasBadSpace(str) {
			str = CleanupBadSpaces(str)
		}
		if HasAdjacentSpaces(str) {
			str = CompressRunsOfSpaces(str)
		}
	}
	if deAccent {
		if ascii {
			str = TransformAccents(str, true, false)
			if HasAdjacentSpaces(str) {
				str = CompressRunsOfSpaces(str)
			}
		}
	}
	if deSymbol {
		if ascii {
			str = FixMisusedLetters(str, false, true, false)
			str = TransformAccents(str, true, false)
			if HasAdjacentSpaces(str) {
				str = CompressRunsOfSpaces(str)
			}
		}
	}
	if doASCII {
		if ascii {
			str = UnicodeToASCII(str)
		}
	}

	if HasFlankingSpace(str) {
		str = strings.TrimSpace(str)
	}

	return str
}

// CleanupPlain removes embedded mixed-content markup tags
func CleanupPlain(str string, wrp bool) string {

	if str == "" {
		return str
	}

	str = strings.Replace(str, "\n", " ", -1)

	str = RemoveEmbeddedMarkup(str)
	str = TransformAccents(str, false, false)

	if HasUnicodeMarkup(str) {
		str = RepairUnicodeMarkup(str, SPACE)
	}
	if HasAngleOrAmpersandEncoding(str) {
		str = RepairTableMarkup(str, SPACE)
		str = RemoveHTMLDecorations(str)
		if wrp {
			str = encodeAngleBracketsAndAmpersand(str)
		}
	}
	if HasBadSpace(str) {
		str = CleanupBadSpaces(str)
	}
	if HasAdjacentSpaces(str) {
		str = CompressRunsOfSpaces(str)
	}
	if HasExtraSpaces(str) {
		str = RemoveExtraSpaces(str)
	}

	return str
}

// CleanupProse converts text to ASCII
func CleanupProse(str string, wrp bool) string {

	if str == "" {
		return str
	}

	str = strings.Replace(str, "\n", " ", -1)

	// CleanupBadSpaces must be called before TransformAccents, to convert Unicode En Space converted from &#x2002;
	if HasBadSpace(str) {
		str = CleanupBadSpaces(str)
	}

	if wrp {
		str = html.UnescapeString(str)
	}

	str = RemoveEmbeddedMarkup(str)
	str = TransformAccents(str, false, false)
	str = FixMisusedLetters(str, true, false, true)
	str = TransformAccents(str, false, false)

	if wrp {
		str = html.EscapeString(str)
	}

	if HasUnicodeMarkup(str) {
		str = RepairUnicodeMarkup(str, SPACE)
	}
	if HasAngleOrAmpersandEncoding(str) {
		str = RepairTableMarkup(str, SPACE)
		str = RemoveHTMLDecorations(str)
		if wrp {
			str = encodeAngleBracketsAndAmpersand(str)
		}
	}
	if HasBadSpace(str) {
		str = CleanupBadSpaces(str)
	}
	if HasAdjacentSpaces(str) {
		str = CompressRunsOfSpaces(str)
	}
	if HasExtraSpaces(str) {
		str = RemoveExtraSpaces(str)
	}

	return str
}

// CleanupQuery performs optional operations on XML query strings
func CleanupQuery(str string, exactMatch, removeBrackets bool) string {

	if exactMatch {
		str = html.EscapeString(str)
	}

	// cleanup string
	if IsNotASCII(str) {
		str = TransformAccents(str, true, false)
		if HasUnicodeMarkup(str) {
			str = RepairUnicodeMarkup(str, SPACE)
		}
	}

	if exactMatch {
		str = strings.ToLower(str)
	}

	if HasBadSpace(str) {
		str = CleanupBadSpaces(str)
	}

	if removeBrackets {
		if HasAngleOrAmpersandEncoding(str) {
			str = RepairEncodedMarkup(str)
			str = RepairScriptMarkup(str, SPACE)
			str = RepairMathMLMarkup(str, SPACE)
			// RemoveEmbeddedMarkup must be called before UnescapeString, which was suppressed in ExploreElements
			str = RemoveEmbeddedMarkup(str)
		}
	}

	if HasAmpOrNotASCII(str) {
		str = html.UnescapeString(str)
	}

	if IsNotASCII(str) {
		if !exactMatch {
			str = UnicodeToASCII(str)
		}
	}
	if HasAdjacentSpaces(str) {
		str = CompressRunsOfSpaces(str)
	}

	return str
}

// CleanupSimple normalizes accented letters, spell Greek letters
func CleanupSimple(str string, wrp bool) string {

	if str == "" {
		return str
	}

	str = strings.Replace(str, "\n", " ", -1)

	str = TransformAccents(str, true, false)

	if HasUnicodeMarkup(str) {
		str = RepairUnicodeMarkup(str, SPACE)
	}
	if HasAngleOrAmpersandEncoding(str) {
		str = RepairTableMarkup(str, SPACE)
		str = RemoveHTMLDecorations(str)
		if wrp {
			str = encodeAngleBracketsAndAmpersand(str)
		}
	}
	if HasBadSpace(str) {
		str = CleanupBadSpaces(str)
	}
	if HasAdjacentSpaces(str) {
		str = CompressRunsOfSpaces(str)
	}
	if HasExtraSpaces(str) {
		str = RemoveExtraSpaces(str)
	}

	return str
}

// CompareAlphaOrNumericKeys is called by slices.Sort{ed}Func to sort strings in alphabetical or numeric order
func CompareAlphaOrNumericKeys(i, j string) int {

	if IsAllDigits(i) && IsAllDigits(j) {
		// numeric sort on strings checks lengths first
		lni, lnj := len(i), len(j)
		// shorter string is numerically lower, assuming no leading zeros
		if lni < lnj {
			return -1
		}
		if lni > lnj {
			return 1
		}
		// same length allows comparison of numeric values without
		// needing expensive string-to-integer conversion
	}

	// same length or non-numeric, can now do string comparison on contents
	return cmp.Compare(i, j)
}

// CompressRunsOfSpaces turns runs of spaces into a single space
func CompressRunsOfSpaces(str string) string {

	// return strings.Join(strings.Fields(str), " ")

	whiteSpace := false
	var buffer strings.Builder

	for _, ch := range str {
		if ch < 127 && inBlank[ch] {
			if !whiteSpace {
				buffer.WriteRune(' ')
			}
			whiteSpace = true
		} else {
			buffer.WriteRune(ch)
			whiteSpace = false
		}
	}

	return buffer.String()
}

// ConvertSlash maps backslash-prefixed letter to ASCII control character
func ConvertSlash(str string) string {

	if str == "" {
		return str
	}

	length := len(str)
	res := make([]byte, length+1, length+1)

	isSlash := false
	idx := 0
	for _, ch := range str {
		if isSlash {
			switch ch {
			case 'n':
				// line feed
				res[idx] = '\n'
			case 'r':
				// carriage return
				res[idx] = '\r'
			case 't':
				// horizontal tab
				res[idx] = '\t'
			case 'f':
				// form feed
				res[idx] = '\f'
			case 'a':
				// audible bell from terminal (undocumented)
				res[idx] = '\x07'
			default:
				res[idx] = byte(ch)
			}
			idx++
			isSlash = false
		} else if ch == '\\' {
			isSlash = true
		} else {
			res[idx] = byte(ch)
			idx++
		}
	}

	res = res[0:idx]

	return string(res)
}

// DoTrimFlankingHTML removes HTML decorations at beginning or end of string
func DoTrimFlankingHTML(str string) string {

	badPrefix := [10]string{
		"<i></i>",
		"<b></b>",
		"<u></u>",
		"<sup></sup>",
		"<sub></sub>",
		"</i>",
		"</b>",
		"</u>",
		"</sup>",
		"</sub>",
	}

	badSuffix := [10]string{
		"<i></i>",
		"<b></b>",
		"<u></u>",
		"<sup></sup>",
		"<sub></sub>",
		"<i>",
		"<b>",
		"<u>",
		"<sup>",
		"<sub>",
	}

	if strings.Contains(str, "<") {
		goOn := true
		for goOn {
			goOn = false
			for _, tag := range badPrefix {
				if strings.HasPrefix(str, tag) {
					str = str[len(tag):]
					goOn = true
				}
			}
			for _, tag := range badSuffix {
				if strings.HasSuffix(str, tag) {
					str = str[:len(str)-len(tag)]
					goOn = true
				}
			}
		}
	}

	return str
}

// FixSpecialCases fixes hyphenated or primed prefixes and suffixes for indexing
func FixSpecialCases(str string) string {

	var arry []string
	var buffer strings.Builder

	terms := strings.Fields(str)

	for _, item := range terms {

		buffer.Reset()

		for i, ch := range item {
			if ch == '-' {
				_, ok := hyphenatedPrefixes[item[0:i]]
				if ok {
					continue
				}
			} else if ch == '\'' {
				_, ok := primedPrefixes[item[0:i]]
				if ok {
					buffer.WriteString("_prime ")
					continue
				}
				_, ok = primedSuffix[item[i:]]
				if ok {
					continue
				}
			}
			buffer.WriteRune(ch)
		}

		item = buffer.String()

		arry = append(arry, item)
	}

	// reconstruct string from transformed words
	str = strings.Join(arry, " ")

	return str
}

// FlattenMathML removes embedded MathML structure
func FlattenMathML(str string, policy int) string {

	findNextXMLBlock := func(txt string) (int, int, bool) {

		beg := strings.Index(txt, "<")
		if beg < 0 {
			return -1, -1, false
		}
		end := strings.Index(txt, ">")
		if end < 0 {
			return -1, -1, false
		}
		end++
		return beg, end, true
	}

	var arry []string

	for {
		beg, end, ok := findNextXMLBlock(str)
		if !ok {
			break
		}
		pfx := str[:beg]
		pfx = strings.TrimSpace(pfx)
		if pfx != "" {
			arry = append(arry, pfx)
		}
		tmp := str[beg:end]
		tmp = strings.TrimSpace(tmp)
		str = str[end:]
	}

	switch policy {
	case PERIOD:
	case SPACE:
		str = strings.Join(arry, " ")
	case CONCISE:
	case BRACKETS:
	case MARKDOWN:
	case SLASH:
	case TAGS:
	case TERSE:
		str = strings.Join(arry, "")
	}

	str = strings.TrimSpace(str)

	// str = RemoveEmbeddedMarkup(str)

	return str
}

// GenBankToMedlineAuthors changes punctuation to make authors searchable
func GenBankToMedlineAuthors(name string) string {

	if name == "" {
		return name
	}

	// start with "Smith-Jones,J.-P."
	// change comma to space
	name = strings.Replace(name, ",", " ", -1)
	// remove all periods
	name = strings.Replace(name, ".", "", -1)
	name = strings.TrimSpace(name)
	idx := strings.LastIndex(name, " ")
	if idx >= 0 {
		// remove hyphen only from initials
		rgt := name[idx:]
		lft := name[:idx]
		rgt = strings.Replace(rgt, "-", "", -1)
		lft = strings.TrimSpace(lft)
		rgt = strings.TrimSpace(rgt)
		name = lft + " " + rgt
	}
	// end with "Smith-Jones JP"

	return name
}

// HasAdjacentSpaces reports if CompressRunsOfSpaces is needed
func HasAdjacentSpaces(str string) bool {

	whiteSpace := false

	for _, ch := range str {
		if ch == ' ' || ch == '\n' {
			if whiteSpace {
				return true
			}
			whiteSpace = true
		} else {
			whiteSpace = false
		}
	}

	return false
}

// HasAdjacentSpacesOrNewline reports if CompressRunsOfSpaces is needed
func HasAdjacentSpacesOrNewline(str string) bool {

	whiteSpace := false

	for _, ch := range str {
		if ch == '\n' {
			return true
		}
		if ch == ' ' {
			if whiteSpace {
				return true
			}
			whiteSpace = true
		} else {
			whiteSpace = false
		}
	}

	return false
}

// HasAmpOrNotASCII reports if UnescapeString is needed
func HasAmpOrNotASCII(str string) bool {

	for _, ch := range str {
		if ch == '&' || ch > 127 {
			return true
		}
	}

	return false
}

// HasAngleOrAmpersandEncoding reports angle brackets or ampersand encodings
func HasAngleOrAmpersandEncoding(str string) bool {

	hasAmp := false
	hasSemi := false

	for _, ch := range str {
		if ch == '<' || ch == '>' {
			return true
		} else if ch == '&' {
			hasAmp = true
		} else if ch == ';' {
			hasSemi = true
		}
	}

	if hasAmp && hasSemi {
		/*
			if strings.Contains(str, "&lt;") ||
				strings.Contains(str, "&gt;") ||
				strings.Contains(str, "&amp;") {
				return true
			}
		*/
		return true
	}

	return false
}

// HasBadSpace matches Unicode spaces other than ASCII space
func HasBadSpace(str string) bool {

	for _, ch := range str {
		if ch > 127 {
			if unicode.IsSpace(ch) {
				return true
			}
			// control characters
			if ch >= 0x0080 && ch <= 0x009F {
				return true
			}
			// no-break space
			if ch == 0x00A0 {
				return true
			}
			// special spaces (en space through zero width space)
			if ch >= 0x2002 && ch <= 0x200B {
				return true
			}
			// narrow no-break space
			if ch == 0x202F {
				return true
			}
			// ideographic space
			if ch == 0x3000 {
				return true
			}
			// unassigned values (private use)
			if ch >= 0xE000 && ch <= 0xF8FF {
				return true
			}
			// zero width no-break space
			if ch == 0xFEFF {
				return true
			}
		}
	}

	return false
}

// HasCombiningAccent reports on spelled-out combining accents
func HasCombiningAccent(str string) bool {

	// will search with components of lower-case combiningAccents map
	str = strings.ToLower(str)

	// do not require whole-phrase context or leading square bracket
	if strings.Index(str, "combining") >= 0 {
		for ky := range lcBadAccents {
			pos := strings.Index(str, ky)
			if pos >= 0 {
				return true
			}
		}
	}

	return false
}

// HasCommaOrSemicolon reports on comma, semicolon, or hyphen
func HasCommaOrSemicolon(str string) bool {

	for _, ch := range str {
		if ch == ',' || ch == ';' || ch == '-' {
			return true
		}
	}

	return false
}

// HasExtraSpaces determines whether RemoveExtraSpaces needs to be called
func HasExtraSpaces(str string) bool {

	if len(str) < 2 {
		return false
	}

	var prev rune

	for _, ch := range str {
		if prev == '(' && ch == ' ' {
			return true
		}
		if prev == ' ' && ch == ')' {
			return true
		}
		if prev == '-' && ch == ' ' {
			return true
		}
		if prev == ' ' && ch == ',' {
			return true
		}
		prev = ch
	}

	return false
}

// HasFlankingSpace returns true if string starts or ends with whitespace
func HasFlankingSpace(str string) bool {

	if str == "" {
		return false
	}

	ch := str[0]
	if ch < 127 && inBlank[ch] {
		return true
	}

	strlen := len(str)
	ch = str[strlen-1]
	if ch < 127 && inBlank[ch] {
		return true
	}

	return false
}

// HasHyphenOrApostrophe reports if string may need indexing preprocessing
func HasHyphenOrApostrophe(str string) bool {

	for _, ch := range str {
		if ch == '-' || ch == '\'' {
			return true
		}
	}

	return false
}

// HasInvisibleUnicode checks for invisible characters
func HasInvisibleUnicode(str string) bool {

	for _, ch := range str {
		if invisibleRunes[ch] {
			return true
		}
	}

	return false
}

// HasLeadingUnicodeDash checks for Unicode dash at beginning of argument
func HasLeadingUnicodeDash(str string) (rune, bool) {

	// intentional Unicode dash at beginning of non-command argument must be escaped by a backslash

	for _, ch := range str {

		// test 32-bit rune
		if ch >= 2010 && ch <= 8213 {
			return ch, true
		}

		if ch == '-' {
			// keep checking for dash after leading hyphen
			continue
		}

		// return once past hyphens
		break
	}

	return 0, false
}

// HasNonASCII checks for non-ASCII Unicode characters
func HasNonASCII(str string) bool {

	for _, ch := range str {
		if ch > 127 {
			return true
		}
	}

	return false
}

// HasSpaceOrHyphen checks for multiple words
func HasSpaceOrHyphen(str string) bool {

	for _, ch := range str {
		if ch == ' ' || ch == '-' {
			return true
		}
	}

	return false
}

// HasUnicodeMarkup checks for Unicode superscript or subscript characters
func HasUnicodeMarkup(str string) bool {

	for _, ch := range str {
		if ch <= 127 {
			continue
		}
		// check for Unicode superscript or subscript characters
		if ch == 0x00B2 || ch == 0x00B3 || ch == 0x00B9 || (ch >= 0x2070 && ch <= 0x208E) {
			return true
		}
	}

	return false
}

// HTMLAhead returns length if next symbol is mixed-content markup tag
func HTMLAhead(text string, idx, txtlen int) int {

	// record position of < character
	start := idx

	// at start of element
	idx++
	if idx >= txtlen {
		return 0
	}
	ch := text[idx]

	if ch == '/' {
		// skip past end tag symbol
		idx++
		ch = text[idx]
	}

	// all embedded markup tags start with a lower-case letter
	if ch < 'a' || ch > 'z' {
		// except for DispFormula in PubmedArticle
		if ch == 'D' && strings.HasPrefix(text[idx:], "DispFormula") {
			for ch != '>' {
				idx++
				ch = text[idx]
			}
			return idx + 1 - start
		}

		// otherwise not a recognized markup tag
		return 0
	}

	idx++
	ch = text[idx]
	for inLower[ch] {
		idx++
		ch = text[idx]
	}

	// if tag name was not all lower-case, then exit
	if ch >= 'A' && ch <= 'Z' {
		return 0
	}

	// skip to end of element, past any attributes or slash character
	for ch != '>' {
		idx++
		ch = text[idx]
	}

	// return number of characters to advance to skip this markup tag
	return idx + 1 - start
}

// HTMLRepair unexpands ampersand encoding to mixed-content formatting symbols
func HTMLRepair(str string) (string, bool) {

	res, ok := htmlRepair[str]

	return res, ok
}

// IsAllCaps matches upper-case
func IsAllCaps(str string) bool {

	for _, ch := range str {
		if !unicode.IsUpper(ch) {
			return false
		}
	}

	return true
}

// IsAllCapsOrDigits matches upper-case or digits
func IsAllCapsOrDigits(str string) bool {

	for _, ch := range str {
		if !unicode.IsUpper(ch) && !unicode.IsDigit(ch) {
			return false
		}
	}

	return true
}

// IsAllDigits matches only digits
func IsAllDigits(str string) bool {

	for _, ch := range str {
		if !unicode.IsDigit(ch) {
			return false
		}
	}

	return true
}

// IsAllDigitsOrPeriod accepts digits or a decimal point
func IsAllDigitsOrPeriod(str string) bool {

	for _, ch := range str {
		if !unicode.IsDigit(ch) && ch != '.' {
			return false
		}
	}

	return true
}

// IsAllLetters matches only letters
func IsAllLetters(str string) bool {

	for _, ch := range str {
		if !unicode.IsLetter(ch) {
			return false
		}
	}

	return true
}

// IsAllLettersOrDigits matches letters or digits
func IsAllLettersOrDigits(str string) bool {

	for _, ch := range str {
		if !unicode.IsLetter(ch) && !unicode.IsDigit(ch) {
			return false
		}
	}

	return true
}

// IsMixedCaps matches if word has both upper- and lower-case letters
// other than initial capital followed by all lower-case
func IsMixedCaps(str string) bool {

	numLower, nonInitialUpper := 0, 0

	for i, ch := range str {
		if unicode.IsLower(ch) {
			numLower++
		}
		if i > 0 && unicode.IsUpper(ch) {
			nonInitialUpper++
		}
	}

	if numLower > 0 && nonInitialUpper > 0 {
		return true
	}

	return false
}

// IsNotASCII returns true for any character greater than 7-bits
func IsNotASCII(str string) bool {

	for _, ch := range str {
		if ch > 127 {
			return true
		}
	}

	return false
}

// IsNotJustWhitespace returns true on any visible or non-whitespace control character
func IsNotJustWhitespace(str string) bool {

	for _, ch := range str {
		if ch > 127 || !inBlank[ch] {
			return true
		}
	}

	return false
}

var plock sync.RWMutex

// IsStopWord returns true for a stop word
func IsStopWord(str string) bool {

	plock.RLock()
	isSW := isStopWord[str]
	plock.RUnlock()

	return isSW
}

// IsUnicodeSubsc matches Unicode subscript characters
func IsUnicodeSubsc(ch rune) bool {
	return ch >= 0x2080 && ch <= 0x209C
}

// IsUnicodeSuper matches Unicode superscript characters
func IsUnicodeSuper(ch rune) bool {
	return ch == 0x00B2 || ch == 0x00B3 || ch == 0x00B9 || (ch >= 0x2070 && ch <= 0x207F)
}

// Ncbi2naToIupac converts a hex-encoded ncbi2na binary nucleotide sequence to IUPAC
func Ncbi2naToIupac(str string) string {

	if str == "" {
		return ""
	}

	var buffer strings.Builder

	src := []byte(str)
	dst := make([]byte, hex.DecodedLen(len(src)))

	n, err := hex.Decode(dst, src)
	if err != nil {
		return ""
	}

	dst = dst[:n]
	for _, byt := range dst {
		tmp := ncbi2naToIupac[int(byt)]
		buffer.WriteString(tmp)
	}

	return buffer.String()
}

// Ncbi4naToIupac converts a hex-encoded ncbi4na binary nucleotide sequence to IUPAC
func Ncbi4naToIupac(str string) string {

	if str == "" {
		return ""
	}

	var buffer strings.Builder

	src := []byte(str)
	dst := make([]byte, hex.DecodedLen(len(src)))

	n, err := hex.Decode(dst, src)
	if err != nil {
		return ""
	}

	dst = dst[:n]
	for _, byt := range dst {
		tmp := ncbi4naToIupac[int(byt)]
		buffer.WriteString(tmp)
	}

	return buffer.String()
}

// NeedsTightening determines whether TightenParentheses needs to be called
func NeedsTightening(str string) bool {

	if len(str) < 2 {
		return false
	}

	var prev rune

	for _, ch := range str {
		if prev == '(' && ch == ' ' {
			return true
		}
		if prev == ' ' && ch == ')' {
			return true
		}
		prev = ch
	}

	return false
}

// NormalizeAuthor is used for citation matching
func NormalizeAuthor(str string) string {

	if str == "" {
		return str
	}

	str = FixMisusedLetters(str, true, true, false)
	str = TransformAccents(str, false, false)

	// convert numeric encoding to apostrophe
	str = strings.Replace(str, "&#39;", "'", -1)
	// remove space following apostrophe
	str = strings.Replace(str, "' ", "'", -1)
	// then remove apostrophe to match indexing logic
	str = strings.Replace(str, "'", "", -1)

	if HasUnicodeMarkup(str) {
		str = RepairUnicodeMarkup(str, SPACE)
	}
	if HasAngleOrAmpersandEncoding(str) {
		str = RepairTableMarkup(str, SPACE)
		// if wrp // str = encodeAngleBracketsAndAmpersand(str)
	}
	if HasBadSpace(str) {
		str = CleanupBadSpaces(str)
	}
	str = strings.Replace(str, "(", " ", -1)
	str = strings.Replace(str, ")", " ", -1)
	if HasAdjacentSpaces(str) {
		str = CompressRunsOfSpaces(str)
	}
	if HasExtraSpaces(str) {
		str = RemoveExtraSpaces(str)
	}

	// remove jr and sr suffix
	if strings.Index(str, " ") != strings.LastIndex(str, " ") {
		str = strings.TrimSuffix(str, " jr")
		str = strings.TrimSuffix(str, " Jr")
		str = strings.TrimSuffix(str, " sr")
		str = strings.TrimSuffix(str, " Sr")
	}

	// normalize remaining punctuation
	str = GenBankToMedlineAuthors(str)

	// truncate to single initial
	pos := strings.LastIndex(str, " ")
	if pos > 0 && len(str) > pos+2 {
		str = str[:pos+2]
	}

	str = CompressRunsOfSpaces(str)
	str = strings.TrimSpace(str)

	return str
}

// NormalizeJournal is used for citation matching
func NormalizeJournal(str string) string {

	if str == "" {
		return str
	}

	// convert numeric encoding to apostrophe
	str = strings.Replace(str, "&#39;", "'", -1)
	// remove apostrophe, but do not remove trailing space
	str = strings.Replace(str, "'", "", -1)

	// convert period to space
	str = strings.Replace(str, ".", " ", -1)

	// convert " & " to " and "
	str = strings.Replace(str, " &amp; ", " and ", -1)

	words := strings.FieldsFunc(str, func(c rune) bool {
		// split at non-alphanumeric characters
		return (!unicode.IsLetter(c) && !unicode.IsDigit(c)) || c > 127
	})

	// stop words to lower case unless first word of journal name or all capital letters or mixed case
	var arry []string

	csr := cases.Title(language.English)

	isFirst := true
	for _, item := range words {

		if !IsAllCaps(item) && !IsMixedCaps(item) {
			item = strings.ToLower(item)
			if !IsStopWord(item) || isFirst {
				item = csr.String(item)
			}
		}

		arry = append(arry, item)
		isFirst = false
	}

	str = strings.Join(arry, " ")

	str = CompressRunsOfSpaces(str)
	str = strings.TrimSpace(str)

	return str
}

// NormalizePage is used for citation matching
func NormalizePage(str string) string {

	if str == "" {
		return str
	}

	// keep starting page before hyphen
	str, _ = SplitInTwoLeft(str, "-")

	str = CompressRunsOfSpaces(str)
	str = strings.TrimSpace(str)

	return str
}

// NormalizeTitle is used for citation matching
func NormalizeTitle(str string) string {

	if str == "" {
		return str
	}

	str = FixMisusedLetters(str, true, false, true)
	str = TransformAccents(str, false, false)
	if HasUnicodeMarkup(str) {
		str = RepairUnicodeMarkup(str, SPACE)
	}
	if HasAngleOrAmpersandEncoding(str) {
		str = RepairTableMarkup(str, SPACE)
		// if wrp // str = encodeAngleBracketsAndAmpersand(str)
	}
	if HasAmpOrNotASCII(str) {
		str = html.UnescapeString(str)
	}
	if HasBadSpace(str) {
		str = CleanupBadSpaces(str)
	}
	if HasAdjacentSpaces(str) {
		str = CompressRunsOfSpaces(str)
	}
	if HasExtraSpaces(str) {
		str = RemoveExtraSpaces(str)
	}

	str = strings.Replace(str, "(", " ", -1)
	str = strings.Replace(str, ")", " ", -1)

	str = strings.Replace(str, "_", " ", -1)

	str = CompressRunsOfSpaces(str)
	str = strings.TrimSpace(str)

	return str
}

// ParseIndex populates the XMLFind structure for matching by identifier
func ParseIndex(indx string) *XMLFind {

	if indx == "" {
		return &XMLFind{}
	}

	// parse parent/element@attribute^version index
	prnt, match := SplitInTwoRight(indx, "/")
	match, versn := SplitInTwoLeft(match, "^")
	match, attrib := SplitInTwoLeft(match, "@")

	return &XMLFind{Index: indx, Parent: prnt, Match: match, Attrib: attrib, Versn: versn}
}

// PrepareForIndexing performs cleanup and normalization of index and query strings
func PrepareForIndexing(str string, doHomoglyphs, isAuthor, isProse, spellGreek, reEncode bool) string {

	if IsNotASCII(str) {
		str = FixMisusedLetters(str, doHomoglyphs, isAuthor, isProse)
		str = TransformAccents(str, spellGreek, reEncode)
		if HasUnicodeMarkup(str) {
			str = RepairUnicodeMarkup(str, SPACE)
		}
	}

	str = strings.ToLower(str)

	if HasBadSpace(str) {
		str = CleanupBadSpaces(str)
	}
	if HasAngleOrAmpersandEncoding(str) {
		str = RepairEncodedMarkup(str)
		str = RepairTableMarkup(str, SPACE)
		str = RepairScriptMarkup(str, SPACE)
		str = RepairMathMLMarkup(str, SPACE)
		// RemoveEmbeddedMarkup must be called before UnescapeString, which was suppressed in ExploreElements
		str = RemoveEmbeddedMarkup(str)
	}

	if HasAmpOrNotASCII(str) {
		str = html.UnescapeString(str)
		str = strings.ToLower(str)
	}

	if HasAdjacentSpaces(str) {
		str = CompressRunsOfSpaces(str)
	}

	str = strings.Replace(str, "(", " ", -1)
	str = strings.Replace(str, ")", " ", -1)

	str = strings.Replace(str, "_", " ", -1)

	if HasHyphenOrApostrophe(str) {
		str = FixSpecialCases(str)
	}

	str = strings.Replace(str, "_", " ", -1)
	str = strings.Replace(str, "-", " ", -1)
	str = strings.Replace(str, "+", " ", -1)
	str = strings.Replace(str, "~", " ", -1)

	// remove trailing punctuation from each word
	var arry []string

	terms := strings.Fields(str)
	for _, item := range terms {
		max := len(item)
		for max > 1 {
			ch := item[max-1]
			if ch != '.' && ch != ',' && ch != ':' && ch != ';' {
				break
			}
			// trim trailing period, comma, colon, and semicolon
			item = item[:max-1]
			// continue checking for runs of punctuation at end
			max--
		}
		if item == "" {
			continue
		}
		arry = append(arry, item)
	}

	// rejoin into string
	str = strings.Join(arry, " ")

	return str
}

// RelaxString removes all punctuation from patterns and search text
func RelaxString(str string) string {

	if str == "" {
		return str
	}

	terms := strings.FieldsFunc(str, func(c rune) bool {
		// split at non-alphanumeric characters
		return (!unicode.IsLetter(c) && !unicode.IsDigit(c)) || c > 127
	})

	str = strings.Join(terms, " ")

	str = CompressRunsOfSpaces(str)
	str = strings.TrimSpace(str)

	return str
}

// RemoveEmbeddedMarkup removes internal mixed-content sections
func RemoveEmbeddedMarkup(str string) string {

	inContent := true
	var buffer strings.Builder

	for _, ch := range str {
		if ch == '<' {
			inContent = false
		} else if ch == '>' {
			inContent = true
		} else if inContent {
			buffer.WriteRune(ch)
		}
	}

	return buffer.String()
}

// RemoveCommaOrSemicolon replaces comma or semicolon with space
func RemoveCommaOrSemicolon(str string) string {

	str = strings.ToLower(str)

	if HasCommaOrSemicolon(str) {
		str = strings.Replace(str, ",", " ", -1)
		str = strings.Replace(str, ";", " ", -1)
		str = CompressRunsOfSpaces(str)
	}
	str = strings.TrimSpace(str)
	str = strings.TrimRight(str, ".?:")

	return str
}

// RemoveExtraSpaces removes spaces inside parentheses, after hyphen, and before comma
func RemoveExtraSpaces(str string) string {

	if len(str) < 2 {
		return str
	}

	var (
		buffer strings.Builder
		prev   rune
	)

	for _, ch := range str {
		if prev == '(' && ch == ' ' {
			ch = '('
		} else if prev == ' ' && ch == ')' {
			ch = ')'
		} else if prev == '-' && ch == ' ' {
			ch = '-'
		} else if prev == ' ' && ch == ',' {
			ch = ','
		} else if prev != 0 {
			buffer.WriteRune(prev)
		}
		prev = ch
	}

	buffer.WriteRune(prev)

	return buffer.String()
}

// RemoveHTMLDecorations removes mixed-content formatting decorations
func RemoveHTMLDecorations(str string) string {

	// unescape converts &lt;b&gt; to <b>
	// also &#181; to Greek mu character
	// also &#x...; where x is followed by a hex value
	str = html.UnescapeString(str)
	// remove mixed content tags
	str = mfix.Replace(str)

	return str
}

// RepairEncodedMarkup removes ampersand-encoded markup
func RepairEncodedMarkup(str string) string {

	// convert &lt;sup&gt; to <sup> (html subset)
	// convert &amp;#181; to &#181; (but not further - use html.UnescapeString)
	// convert &amp;amp;amp;amp;amp;amp;amp;lt; to &lt;
	// remove </sub><sub> or </sup><sup> (internals)

	var buffer strings.Builder

	lookAhead := func(txt string, to int) string {
		mx := len(txt)
		if to > mx {
			to = mx
		}
		pos := strings.Index(txt[:to], "gt;")
		if pos > 0 {
			to = pos + 3
		}
		return txt[:to]
	}

	skip := 0

	for i, ch := range str {
		if skip > 0 {
			skip--
			continue
		}
		if ch == '<' {
			// remove internal tags in runs of subscripts or superscripts
			if strings.HasPrefix(str[i:], "</sub><sub>") || strings.HasPrefix(str[i:], "</sup><sup>") {
				skip = 10
				continue
			}
			buffer.WriteRune(ch)
			continue
		} else if ch != '&' {
			buffer.WriteRune(ch)
			continue
		} else if strings.HasPrefix(str[i:], "&lt;") {
			sub := lookAhead(str[i:], 14)
			txt, ok := htmlRepair[sub]
			if ok {
				adv := len(sub) - 1
				// do not convert if flanked by spaces - it may be a scientific symbol,
				// e.g., fragments <i> in PMID 9698410, or escaped <b> and <d> tags used
				// to indicate stem position in letters in PMID 21892341
				if i < 1 || str[i-1] != ' ' || !strings.HasPrefix(str[i+adv:], "; ") {
					buffer.WriteString(txt)
					skip = adv
					continue
				}
			}
		} else if strings.HasPrefix(str[i:], "&amp;") {
			if strings.HasPrefix(str[i:], "&amp;lt;") {
				sub := lookAhead(str[i:], 22)
				txt, ok := htmlRepair[sub]
				if ok {
					buffer.WriteString(txt)
					skip = len(sub) - 1
					continue
				} else {
					buffer.WriteString("&lt;")
					skip = 7
					continue
				}
			} else if strings.HasPrefix(str[i:], "&amp;gt;") {
				buffer.WriteString("&gt;")
				skip = 7
				continue
			} else {
				skip = 4
				j := i + 5
				// remove runs of multiply-encoded ampersands
				for strings.HasPrefix(str[j:], "amp;") {
					skip += 4
					j += 4
				}
				// then look for special symbols used in PubMed records
				if strings.HasPrefix(str[j:], "lt;") {
					buffer.WriteString("&lt;")
					skip += 3
				} else if strings.HasPrefix(str[j:], "gt;") {
					buffer.WriteString("&gt;")
					skip += 3
				} else if strings.HasPrefix(str[j:], "frac") {
					buffer.WriteString("&frac")
					skip += 4
				} else if strings.HasPrefix(str[j:], "plusmn") {
					buffer.WriteString("&plusmn")
					skip += 6
				} else if strings.HasPrefix(str[j:], "acute") {
					buffer.WriteString("&acute")
					skip += 5
				} else if strings.HasPrefix(str[j:], "aacute") {
					buffer.WriteString("&aacute")
					skip += 6
				} else if strings.HasPrefix(str[j:], "rsquo") {
					buffer.WriteString("&rsquo")
					skip += 5
				} else if strings.HasPrefix(str[j:], "lsquo") {
					buffer.WriteString("&lsquo")
					skip += 5
				} else if strings.HasPrefix(str[j:], "micro") {
					buffer.WriteString("&micro")
					skip += 5
				} else if strings.HasPrefix(str[j:], "oslash") {
					buffer.WriteString("&oslash")
					skip += 6
				} else if strings.HasPrefix(str[j:], "kgr") {
					buffer.WriteString("&kgr")
					skip += 3
				} else if strings.HasPrefix(str[j:], "apos") {
					buffer.WriteString("&apos")
					skip += 4
				} else if strings.HasPrefix(str[j:], "quot") {
					buffer.WriteString("&quot")
					skip += 4
				} else if strings.HasPrefix(str[j:], "alpha") {
					buffer.WriteString("&alpha")
					skip += 5
				} else if strings.HasPrefix(str[j:], "beta") {
					buffer.WriteString("&beta")
					skip += 4
				} else if strings.HasPrefix(str[j:], "gamma") {
					buffer.WriteString("&gamma")
					skip += 5
				} else if strings.HasPrefix(str[j:], "Delta") {
					buffer.WriteString("&Delta")
					skip += 5
				} else if strings.HasPrefix(str[j:], "phi") {
					buffer.WriteString("&phi")
					skip += 3
				} else if strings.HasPrefix(str[j:], "ge") {
					buffer.WriteString("&ge")
					skip += 2
				} else if strings.HasPrefix(str[j:], "sup2") {
					buffer.WriteString("&sup2")
					skip += 4
				} else if strings.HasPrefix(str[j:], "#") {
					buffer.WriteString("&")
				} else {
					buffer.WriteString("&amp;")
				}
				continue
			}
		}

		// if loop not continued by any preceding test, print character
		buffer.WriteRune(ch)
	}

	return buffer.String()
}

// RepairMathMLMarkup removes MathML embedded markup symbols
func RepairMathMLMarkup(str string, policy int) string {

	str = strings.Replace(str, "> <mml:", "><mml:", -1)
	str = strings.Replace(str, "> </mml:", "></mml:", -1)

	findNextMathBlock := func(txt string) (int, int, bool) {

		beg := strings.Index(txt, "<DispFormula")
		if beg < 0 {
			return -1, -1, false
		}
		end := strings.Index(txt, "</DispFormula>")
		if end < 0 {
			return -1, -1, false
		}
		end += 14
		return beg, end, true
	}

	var arry []string

	for {
		beg, end, ok := findNextMathBlock(str)
		if !ok {
			break
		}
		pfx := str[:beg]
		pfx = strings.TrimSpace(pfx)
		arry = append(arry, pfx)
		tmp := str[beg:end]
		if strings.HasPrefix(tmp, "<DispFormula") {
			tmp = FlattenMathML(tmp, policy)
		}
		tmp = strings.TrimSpace(tmp)
		arry = append(arry, tmp)
		str = str[end:]
	}

	str = strings.TrimSpace(str)
	arry = append(arry, str)

	return strings.Join(arry, " ")
}

// RepairScriptMarkup converts mixed-content subscript and superscript tags
func RepairScriptMarkup(str string, policy int) string {

	var buffer strings.Builder

	skip := 0

	for i, ch := range str {
		if skip > 0 {
			skip--
			continue
		}
		if ch == '<' {
			if strings.HasPrefix(str[i:], "<sub>") {
				switch policy {
				case PERIOD:
				case SPACE:
				case CONCISE:
				case BRACKETS:
					buffer.WriteRune('(')
				case MARKDOWN:
					buffer.WriteRune('~')
				}
				skip = 4
				continue
			}
			if strings.HasPrefix(str[i:], "<sup>") {
				switch policy {
				case PERIOD:
				case SPACE:
				case CONCISE:
					buffer.WriteRune('[')
				case BRACKETS:
					buffer.WriteRune('[')
				case MARKDOWN:
					buffer.WriteRune('^')
				}
				skip = 4
				continue
			}
			if strings.HasPrefix(str[i:], "</sub>") {
				if strings.HasPrefix(str[i+6:], "<sup>") {
					switch policy {
					case PERIOD:
						buffer.WriteRune('.')
					case SPACE:
						buffer.WriteRune(' ')
					case CONCISE:
						buffer.WriteRune('[')
					case BRACKETS:
						buffer.WriteRune(')')
						buffer.WriteRune('[')
					case MARKDOWN:
						buffer.WriteRune('~')
						buffer.WriteRune('^')
					}
					skip = 10
					continue
				}
				switch policy {
				case PERIOD:
				case SPACE:
				case CONCISE:
				case BRACKETS:
					buffer.WriteRune(')')
				case MARKDOWN:
					buffer.WriteRune('~')
				}
				skip = 5
				continue
			}
			if strings.HasPrefix(str[i:], "</sup>") {
				if strings.HasPrefix(str[i+6:], "<sub>") {
					switch policy {
					case PERIOD:
						buffer.WriteRune('.')
					case SPACE:
						buffer.WriteRune(' ')
					case CONCISE:
						buffer.WriteRune(']')
					case BRACKETS:
						buffer.WriteRune(']')
						buffer.WriteRune('(')
					case MARKDOWN:
						buffer.WriteRune('^')
						buffer.WriteRune('~')
					}
					skip = 10
					continue
				}
				switch policy {
				case PERIOD:
				case SPACE:
				case CONCISE:
					buffer.WriteRune(']')
				case BRACKETS:
					buffer.WriteRune(']')
				case MARKDOWN:
					buffer.WriteRune('^')
				}
				skip = 5
				continue
			}
		}

		buffer.WriteRune(ch)
	}

	return buffer.String()
}

// RepairTableMarkup removes HTML table element tags
func RepairTableMarkup(str string, policy int) string {

	str = strings.Replace(str, "<tr>", " ", -1)
	str = strings.Replace(str, "<td>", " ", -1)
	str = strings.Replace(str, "</tr>", " ", -1)
	str = strings.Replace(str, "</td>", " ", -1)

	return str
}

// RepairUnicodeMarkup handles Unicode superscript and subscript symbols
func RepairUnicodeMarkup(str string, policy int) string {

	type MarkupType int

	const (
		NOSCRIPT MarkupType = iota
		SUPSCRIPT
		SUBSCRIPT
		PLAINDIGIT
	)

	var buffer strings.Builder

	// to improve readability, keep track of switches between numeric types, add period at transitions when converting to plain ASCII
	level := NOSCRIPT

	for _, ch := range str {
		if ch > 127 {
			rn, ok := superScriptRunes[ch]
			if ok {
				ch = rn
				switch level {
				case NOSCRIPT:
					switch policy {
					case PERIOD:
					case SPACE:
					case CONCISE:
						buffer.WriteRune('[')
					case BRACKETS:
						buffer.WriteRune('[')
					case MARKDOWN:
						buffer.WriteRune('^')
					case SLASH:
					case TAGS:
						buffer.WriteString("<sup>")
					}
				case SUPSCRIPT:
					switch policy {
					case PERIOD:
					case SPACE:
					case CONCISE:
					case BRACKETS:
					case MARKDOWN:
					case SLASH:
					case TAGS:
					}
				case SUBSCRIPT:
					switch policy {
					case PERIOD:
						buffer.WriteRune('.')
					case SPACE:
						buffer.WriteRune(' ')
					case CONCISE:
						buffer.WriteRune('[')
					case BRACKETS:
						buffer.WriteRune(')')
						buffer.WriteRune('[')
					case MARKDOWN:
						buffer.WriteRune('~')
						buffer.WriteRune('^')
					case SLASH:
						buffer.WriteRune('\\')
					case TAGS:
						buffer.WriteString("</sub>")
						buffer.WriteString("<sup>")
					}
				case PLAINDIGIT:
					switch policy {
					case PERIOD:
						buffer.WriteRune('.')
					case SPACE:
						buffer.WriteRune(' ')
					case CONCISE:
						buffer.WriteRune('[')
					case BRACKETS:
						buffer.WriteRune('[')
					case MARKDOWN:
						buffer.WriteRune('^')
					case SLASH:
						buffer.WriteRune('\\')
					case TAGS:
						buffer.WriteString("<sup>")
					}
				}
				level = SUPSCRIPT
				buffer.WriteRune(ch)
				continue
			}
			rn, ok = subScriptRunes[ch]
			if ok {
				ch = rn
				switch level {
				case NOSCRIPT:
					switch policy {
					case PERIOD:
					case SPACE:
					case CONCISE:
					case BRACKETS:
						buffer.WriteRune('(')
					case MARKDOWN:
						buffer.WriteRune('~')
					case SLASH:
					case TAGS:
						buffer.WriteString("<sub>")
					}
				case SUPSCRIPT:
					switch policy {
					case PERIOD:
						buffer.WriteRune('.')
					case SPACE:
						buffer.WriteRune(' ')
					case CONCISE:
						buffer.WriteRune(']')
					case BRACKETS:
						buffer.WriteRune(']')
						buffer.WriteRune('(')
					case MARKDOWN:
						buffer.WriteRune('^')
						buffer.WriteRune('~')
					case SLASH:
						buffer.WriteRune('/')
					case TAGS:
						buffer.WriteString("</sup>")
						buffer.WriteString("<sub>")
					}
				case SUBSCRIPT:
					switch policy {
					case PERIOD:
					case SPACE:
					case CONCISE:
					case BRACKETS:
					case MARKDOWN:
					case SLASH:
					case TAGS:
					}
				case PLAINDIGIT:
					switch policy {
					case PERIOD:
						buffer.WriteRune('.')
					case SPACE:
						buffer.WriteRune(' ')
					case CONCISE:
					case BRACKETS:
						buffer.WriteRune('(')
					case MARKDOWN:
						buffer.WriteRune('~')
					case SLASH:
						buffer.WriteRune('/')
					case TAGS:
						buffer.WriteString("<sub>")
					}
				}
				level = SUBSCRIPT
				buffer.WriteRune(ch)
				continue
			}
			switch level {
			case NOSCRIPT:
				switch policy {
				case PERIOD:
				case SPACE:
				case CONCISE:
				case BRACKETS:
				case MARKDOWN:
				case SLASH:
				case TAGS:
				}
			case SUPSCRIPT:
				switch policy {
				case PERIOD:
				case SPACE:
				case CONCISE:
					buffer.WriteRune(']')
				case BRACKETS:
					buffer.WriteRune(']')
				case MARKDOWN:
					buffer.WriteRune('^')
				case SLASH:
				case TAGS:
					buffer.WriteString("</sup>")
				}
			case SUBSCRIPT:
				switch policy {
				case PERIOD:
				case SPACE:
				case CONCISE:
				case BRACKETS:
					buffer.WriteRune(')')
				case MARKDOWN:
					buffer.WriteRune('~')
				case SLASH:
				case TAGS:
					buffer.WriteString("</sub>")
				}
			case PLAINDIGIT:
				switch policy {
				case PERIOD:
				case SPACE:
				case CONCISE:
				case BRACKETS:
				case MARKDOWN:
				case SLASH:
				case TAGS:
				}
			}
			level = NOSCRIPT
		} else if ch >= '0' && ch <= '9' {
			switch level {
			case NOSCRIPT:
				switch policy {
				case PERIOD:
				case SPACE:
				case CONCISE:
				case BRACKETS:
				case MARKDOWN:
				case SLASH:
				case TAGS:
				}
			case SUPSCRIPT:
				switch policy {
				case PERIOD:
					buffer.WriteRune('.')
				case SPACE:
					buffer.WriteRune(' ')
				case CONCISE:
					buffer.WriteRune(']')
				case BRACKETS:
					buffer.WriteRune(']')
				case MARKDOWN:
					buffer.WriteRune('^')
				case SLASH:
					buffer.WriteRune('/')
				case TAGS:
					buffer.WriteString("</sup>")
				}
			case SUBSCRIPT:
				switch policy {
				case PERIOD:
					buffer.WriteRune('.')
				case SPACE:
					buffer.WriteRune(' ')
				case CONCISE:
				case BRACKETS:
					buffer.WriteRune(')')
				case MARKDOWN:
					buffer.WriteRune('~')
				case SLASH:
					buffer.WriteRune('\\')
				case TAGS:
					buffer.WriteString("</sub>")
				}
			case PLAINDIGIT:
				switch policy {
				case PERIOD:
				case SPACE:
				case CONCISE:
				case BRACKETS:
				case MARKDOWN:
				case SLASH:
				case TAGS:
				}
			}
			level = PLAINDIGIT
		} else {
			switch level {
			case NOSCRIPT:
				switch policy {
				case PERIOD:
				case SPACE:
				case CONCISE:
				case BRACKETS:
				case MARKDOWN:
				case SLASH:
				case TAGS:
				}
			case SUPSCRIPT:
				switch policy {
				case PERIOD:
				case SPACE:
				case CONCISE:
					buffer.WriteRune(']')
				case BRACKETS:
					buffer.WriteRune(']')
				case MARKDOWN:
					buffer.WriteRune('^')
				case SLASH:
				case TAGS:
					buffer.WriteString("</sup>")
				}
			case SUBSCRIPT:
				switch policy {
				case PERIOD:
				case SPACE:
				case CONCISE:
				case BRACKETS:
					buffer.WriteRune(')')
				case MARKDOWN:
					buffer.WriteRune('~')
				case SLASH:
				case TAGS:
					buffer.WriteString("</sub>")
				}
			case PLAINDIGIT:
				switch policy {
				case PERIOD:
				case SPACE:
				case CONCISE:
				case BRACKETS:
				case MARKDOWN:
				case SLASH:
				case TAGS:
				}
			}
			level = NOSCRIPT
		}
		buffer.WriteRune(ch)
	}

	switch level {
	case NOSCRIPT:
		switch policy {
		case PERIOD:
		case SPACE:
		case CONCISE:
		case BRACKETS:
		case MARKDOWN:
		case SLASH:
		case TAGS:
		}
	case SUPSCRIPT:
		switch policy {
		case PERIOD:
		case SPACE:
		case CONCISE:
			buffer.WriteRune(']')
		case BRACKETS:
			buffer.WriteRune(']')
		case MARKDOWN:
			buffer.WriteRune('^')
		case SLASH:
		case TAGS:
			buffer.WriteString("</sup>")
		}
	case SUBSCRIPT:
		switch policy {
		case PERIOD:
		case SPACE:
		case CONCISE:
		case BRACKETS:
			buffer.WriteRune(')')
		case MARKDOWN:
			buffer.WriteRune('~')
		case SLASH:
		case TAGS:
			buffer.WriteString("</sub>")
		}
	case PLAINDIGIT:
		switch policy {
		case PERIOD:
		case SPACE:
		case CONCISE:
		case BRACKETS:
		case MARKDOWN:
		case SLASH:
		case TAGS:
		}
	}

	return buffer.String()
}

// ReverseString creates a mirror-image of letters in a string
func ReverseString(str string) string {

	if str == "" {
		return str
	}

	runes := []rune(str)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	rev := string(runes)

	return rev
}

// SlidingSlices passes a sliding window over a sequence
func SlidingSlices(str string, num int) []string {

	if str == "" {
		return nil
	}

	var arry []string

	fst, lst, max := 0, num, len(str)+1

	for lst < max {
		sub := str[fst:lst]
		arry = append(arry, sub)
		fst++
		lst++
	}

	return arry
}

// SortStringByWords sorts the individual words in a string
func SortStringByWords(str string) string {

	str = RemoveCommaOrSemicolon(str)

	// check for multiple words
	if HasSpaceOrHyphen(str) {
		flds := strings.Fields(str)
		slices.Sort(flds)
		str = strings.Join(flds, " ")
		str = strings.Replace(str, "-", " ", -1)
		str = CompressRunsOfSpaces(str)
		str = strings.TrimRight(str, ".?:")
	}

	return str
}

// SplitInTwoLeft loads the first argument if no delimiter is present
func SplitInTwoLeft(str, chr string) (string, string) {

	slash := strings.SplitN(str, chr, 2)
	if len(slash) > 1 {
		return slash[0], slash[1]
	}

	return str, ""
}

// SplitInTwoRight loads the second argument if no delimiter is present
func SplitInTwoRight(str, chr string) (string, string) {

	slash := strings.SplitN(str, chr, 2)
	if len(slash) > 1 {
		return slash[0], slash[1]
	}

	return "", str
}

// TightenParentheses removes spaces inside parentheses
func TightenParentheses(str string) string {

	if len(str) < 2 {
		return str
	}

	var (
		buffer strings.Builder
		prev   rune
	)

	for _, ch := range str {
		if prev == '(' && ch == ' ' {
			ch = '('
		} else if prev == ' ' && ch == ')' {
			ch = ')'
		} else if prev != 0 {
			buffer.WriteRune(prev)
		}
		prev = ch
	}

	buffer.WriteRune(prev)

	return buffer.String()
}

// TrimBOM removes leading Byte Order Mark, if present
func TrimBOM(str string) string {

	if str == "" {
		return str
	}

	for i, ch := range str {

		if i == 0 {
			if ch == 65279 {
				// next offset is size of BOM
				continue
			}
			// otherwise break out of loop
			break
		}

		// trim with size of BOM (should be 3)
		str = str[i:]

		// break to return trimmed string
		break
	}

	return str
}

// UnicodeToASCII converts nonASCII to ampersand-pound-x encoding
func UnicodeToASCII(str string) string {

	// convert to &#x...; representation

	var buffer strings.Builder

	for _, ch := range str {
		if ch > 127 {
			s := strconv.QuoteToASCII(string(ch))
			s = strings.ToUpper(s[3:7])
			for {
				if !strings.HasPrefix(s, "0") {
					break
				}
				s = s[1:]
			}
			buffer.WriteString("&#x")
			buffer.WriteString(s)
			buffer.WriteRune(';')
			continue
		}
		buffer.WriteRune(ch)
	}

	return buffer.String()
}

// initialize mfix replacer before non-init functions are called
func init() {

	mfix = strings.NewReplacer(
		"<b>", "",
		"<i>", "",
		"<u>", "",
		"</b>", "",
		"</i>", "",
		"</u>", "",
		"<b/>", "",
		"<i/>", "",
		"<u/>", "",
		"<sup>", "",
		"<sub>", "",
		"</sup>", "",
		"</sub>", "",
		"<sup/>", "",
		"<sub/>", "",
	)

	lcBadAccents = make(map[string]bool)

	for ky := range combiningAccents {
		ky = strings.ToLower(ky)
		lcBadAccents[ky] = true
	}

}
