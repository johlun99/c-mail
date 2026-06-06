package mail

// This file holds the illustrative mock dataset, ported from the design
// prototype (design_handoff_cmail/prototype/cmail-data.js). All copy is Swedish.
// Replace with real Gmail data behind the Store interface.

// defaultCategories is the built-in triage taxonomy. Users may add more at
// runtime; these are the defaults.
func defaultCategories() []Category {
	return []Category{
		{Key: "svara", Label: "att svara på", Short: "svara", Color: "c-accent", Desc: "Väntar på ditt svar"},
		{Key: "intressant", Label: "intressant", Short: "intresse", Color: "c-blue", Desc: "Värt att läsa"},
		{Key: "faktura", Label: "faktura", Short: "faktura", Color: "c-red", Desc: "Att betala / förfaller"},
		{Key: "kvitto", Label: "kvitton", Short: "kvitto", Color: "c-amber", Desc: "Köp & bekräftelser"},
		{Key: "reklam", Label: "reklam", Short: "reklam", Color: "c-dim", Desc: "Nyhetsbrev & utskick"},
		{Key: "vantar", Label: "väntar på svar", Short: "väntar", Color: "c-violet", Desc: "Du väntar på dem"},
	}
}

// added builds a slice of agent-authored draft lines from plain strings.
func added(lines ...string) []DraftLine {
	out := make([]DraftLine, len(lines))
	for i, t := range lines {
		out[i] = DraftLine{Text: t, Kind: DraftAdded}
	}
	return out
}

func mockMails() []Mail {
	return []Mail{
		{
			ID: "m1", Cat: "svara", From: "Anna Lindqvist", FromAddr: "anna.lindqvist@nordveda.se",
			Avatar: "AL", Subject: "Re: Offert — Q3 integrationsprojekt",
			Snippet: "Hej! Tack för uppdaterad offert. Vi vill köra igång men har två frågor om tidsplanen innan vi …",
			Time:    "09:12", Day: "idag", Unread: true, Deadline: "fre 6 jun", Confidence: 0.94, ThreadCount: 12,
			Body: []string{
				"Hej!",
				"Tack för den uppdaterade offerten — den ser bra ut och vi vill gärna köra igång. Innan vi skriver på har vi två saker vi behöver reda ut:",
				"1) Tidsplanen säger leverans v.34, men vår release är låst till v.33. Finns det utrymme att tidigarelägga sista milstolpen en vecka?",
				"2) Ingår löpande support de första 30 dagarna efter driftsättning, eller är det en separat post?",
				"Kan vi få svar innan fredag så hinner vi ta beslut på måndagens styrgruppsmöte.",
				"Tack på förhand!",
				"Vänliga hälsningar,\nAnna Lindqvist\nNordveda AB",
			},
			Agent: AgentAnalysis{
				Summary: "Kund vill skriva på men behöver svar på 2 frågor (tidsplan + support) före fre.",
				Facts: []Fact{
					{"kategori", "[att svara på] · 0.94"},
					{"avsändare", "känd · 12 tidigare trådar"},
					{"deadline", "fre 6 jun (styrgrupp må)"},
					{"sentiment", "positiv · köpintention hög"},
					{"åtgärd", "utkast genererat ↓"},
				},
				Tasks: []string{"Bekräfta v.33-leverans", "Klargör 30-dagars support"},
			},
			Draft: &Draft{
				Tone: "professionell · varm",
				Lines: added(
					"Hej Anna,",
					"",
					"Vad roligt att ni vill köra igång! Här kommer svar på era frågor:",
					"",
					"1) Tidsplan: vi kan tidigarelägga sista milstolpen till v.33 — det fungerar bra för oss och vi justerar planen därefter.",
					"2) Support: de första 30 dagarna efter driftsättning ingår utan extra kostnad. Därefter går vi över till löpande supportavtal.",
					"",
					"Säg till om ni vill att jag skickar över ett uppdaterat avtal redan idag.",
					"",
					"Vänliga hälsningar,",
					"Johan",
				),
			},
		},
		{
			ID: "m2", Cat: "svara", From: "David Ek", FromAddr: "david@frostbyte.io",
			Avatar: "DE", Subject: "Snabb fråga om API-nycklarna",
			Snippet: "Har du möjlighet att rotera nyckeln för staging? Den gamla läckte i en gammal commit …",
			Time:    "08:40", Day: "idag", Unread: true, Confidence: 0.88, ThreadCount: 3,
			Body: []string{
				"Tjena,",
				"Vi hittade den gamla staging-nyckeln i en commit från i våras. Inget akut men kan du rotera den och skicka den nya via 1Password?",
				"Tack!",
				"/David",
			},
			Agent: AgentAnalysis{
				Summary: "Begär nyckelrotation för staging. Säkerhetsärende, låg brådska.",
				Facts: []Fact{
					{"kategori", "[att svara på] · 0.88"},
					{"avsändare", "känd · team"},
					{"typ", "åtgärd → uppgift"},
					{"åtgärd", "utkast genererat ↓"},
				},
				Tasks: []string{"Rotera staging-nyckel", "Dela via 1Password"},
			},
			Draft: &Draft{
				Tone: "kort · kollegialt",
				Lines: added(
					"Tjena David,",
					"",
					"Fixar det idag — roterar staging-nyckeln och lägger den nya i 1Password (valvet \"Infra\"). Pingar dig när den är på plats.",
					"",
					"/Johan",
				),
			},
		},
		{
			ID: "m3", Cat: "faktura", From: "Fortnox", FromAddr: "no-reply@fortnox.se",
			Avatar: "FX", Subject: "Faktura 2026-0418 förfaller om 3 dagar",
			Snippet: "Belopp 14 500 kr. Förfallodatum 2026-06-08. OCR 9012 3456 78 …",
			Time:    "07:55", Day: "idag", Unread: true, Deadline: "8 jun", Confidence: 0.99, ThreadCount: 1,
			Body: []string{
				"Detta är en påminnelse om att faktura 2026-0418 förfaller inom kort.",
				"Belopp: 14 500 kr\nFörfallodatum: 2026-06-08\nOCR: 9012 3456 78\nBankgiro: 123-4567",
				"Logga in i Fortnox för att se underlaget.",
			},
			Agent: AgentAnalysis{
				Summary: "Leverantörsfaktura 14 500 kr, förfaller 8 jun. Återkommande (månadsvis SaaS).",
				Facts: []Fact{
					{"kategori", "[faktura] · 0.99"},
					{"belopp", "14 500 kr"},
					{"förfaller", "8 jun · om 3 dagar"},
					{"mönster", "återkommande · 11 tidigare"},
				},
				Tasks: []string{"Lägg till i betalkö", "Matcha mot underlag"},
			},
		},
		{
			ID: "m4", Cat: "intressant", From: "Stratechery", FromAddr: "ben@stratechery.com",
			Avatar: "ST", Subject: "The Agentic Inbox and the End of Triage",
			Snippet: "When the marginal cost of reading drops to zero, the bottleneck moves from attention to …",
			Time:    "06:30", Day: "idag", Unread: false, Confidence: 0.81, ThreadCount: 1,
			Body: []string{
				"When the marginal cost of reading an email drops to zero, the bottleneck moves from attention to trust.",
				"The interesting design question for agentic mail clients is not summarization — it is delegation: how much are you willing to let an agent decide on your behalf, and how do you make that boundary legible?",
				"(Full essay on the web.)",
			},
			Agent: AgentAnalysis{
				Summary: "Essä om agentiska mailklienter & delegering. Matchar dina sparade ämnen.",
				Facts: []Fact{
					{"kategori", "[intressant] · 0.81"},
					{"matchar", "ämne: \"agents\", \"produktivitet\""},
					{"lästid", "~6 min"},
				},
				Tasks: []string{"Spara till läslista"},
			},
		},
		{
			ID: "m5", Cat: "kvitto", From: "Apple", FromAddr: "no_reply@email.apple.com",
			Avatar: "AP", Subject: "Ditt kvitto från App Store",
			Snippet: "Tack för ditt köp. Things 3 — 149,00 kr. Beställnings-ID MX0J2K …",
			Time:    "igår", Day: "igår", Unread: false, Confidence: 0.97, ThreadCount: 1,
			Body: []string{
				"Tack för ditt köp.",
				"Things 3 (Mac) — 149,00 kr\nBeställnings-ID: MX0J2K8L\nDatum: 2026-06-04",
			},
			Agent: AgentAnalysis{
				Summary: "Kvitto, 149 kr, mjukvara. Bokförs som programvara.",
				Facts: []Fact{
					{"kategori", "[kvitto] · 0.97"},
					{"belopp", "149,00 kr"},
					{"konto", "→ 5420 programvara"},
				},
				Tasks: []string{"Exportera till bokföring"},
			},
		},
		{
			ID: "m6", Cat: "vantar", From: "Mira Holm", FromAddr: "mira@tegelvik.se",
			Avatar: "MH", Subject: "Re: Designgranskning — väntar på din feedback",
			Snippet: "Ingen stress, men hörde inget om mockuparna jag skickade i måndags. Hinner du titta …",
			Time:    "igår", Day: "igår", Unread: false, Confidence: 0.86, ThreadCount: 5,
			Body: []string{
				"Hej igen!",
				"Ingen stress, men jag skickade mockuparna i måndags och ville bara stämma av — hinner du titta innan veckans slut?",
				"/Mira",
			},
			Agent: AgentAnalysis{
				Summary: "Du har inte svarat på mockuparna sedan måndag. Mira följer upp.",
				Facts: []Fact{
					{"kategori", "[väntar på svar] · 0.86"},
					{"din skuld", "2 dagars fördröjning"},
					{"åtgärd", "utkast genererat ↓"},
				},
				Tasks: []string{"Granska mockuparna", "Svara Mira"},
			},
			Draft: &Draft{
				Tone: "vänligt · ursäktande",
				Lines: added(
					"Hej Mira,",
					"",
					"Ursäkta dröjsmålet! Jag tittar igenom mockuparna idag och återkommer med samlad feedback senast imorgon förmiddag.",
					"",
					"/Johan",
				),
			},
		},
		{
			ID: "m7", Cat: "reklam", From: "Figma", FromAddr: "news@figma.com",
			Avatar: "FG", Subject: "Nytt i Figma: Config 2026-höjdpunkter",
			Snippet: "Se alla nyheter från årets Config — nya AI-verktyg, dev mode-uppdateringar och mer …",
			Time:    "igår", Day: "igår", Unread: false, Confidence: 0.92, ThreadCount: 1,
			Body: []string{"Nyhetsbrev från Figma. Höjdpunkter från Config 2026."},
			Agent: AgentAnalysis{
				Summary: "Produktnyhetsbrev. Auto-arkiveras enligt din regel för @figma.com.",
				Facts: []Fact{
					{"kategori", "[reklam] · 0.92"},
					{"regel", "auto-arkivera nyhetsbrev"},
				},
				Tasks: []string{},
			},
		},
		{
			ID: "m8", Cat: "svara", From: "Rekrytering · Voltera", FromAddr: "talent@voltera.tech",
			Avatar: "VT", Subject: "Spännande roll som Staff Engineer?",
			Snippet: "Vi bygger nästa generations energilager och letar efter någon med din profil …",
			Time:    "igår", Day: "igår", Unread: false, Confidence: 0.72, ThreadCount: 1,
			Body: []string{
				"Hej Johan,",
				"Vi på Voltera bygger nästa generations energilager och tror att du skulle passa in. Har du 20 min för ett förutsättningslöst samtal nästa vecka?",
				"Vänligen, Sara — Talent, Voltera",
			},
			Agent: AgentAnalysis{
				Summary: "Rekryterare, kall kontakt. Låg konfidens — kan vara [reklam] eller [svara].",
				Facts: []Fact{
					{"kategori", "[att svara på] · 0.72 ⚠"},
					{"avsändare", "okänd"},
					{"not", "låg konfidens — be om bekräftelse"},
					{"åtgärd", "artigt nej-utkast ↓"},
				},
				Tasks: []string{"Bestäm: svara eller arkivera"},
			},
			Draft: &Draft{
				Tone: "artigt · kort nej",
				Lines: added(
					"Hej Sara,",
					"",
					"Tack för att du hörde av dig — kul roll, men jag är inte tillgänglig just nu. Återkom gärna framåt hösten.",
					"",
					"Vänligen, Johan",
				),
			},
		},
		{
			ID: "m9", Cat: "kvitto", From: "SJ", FromAddr: "biljett@sj.se",
			Avatar: "SJ", Subject: "Din biljett: Stockholm → Göteborg",
			Snippet: "Avgång 11 jun 08:21, vagn 4 plats 52. Boknings-nr ZK19QP …",
			Time:    "ons", Day: "tidigare", Unread: false, Confidence: 0.95, ThreadCount: 1,
			Body: []string{"Stockholm C → Göteborg C\n11 jun 08:21–11:34\nVagn 4, plats 52\nBoknings-nr: ZK19QP"},
			Agent: AgentAnalysis{
				Summary: "Tågbiljett 11 jun. Lagt till i din kalender (förslag).",
				Facts: []Fact{
					{"kategori", "[kvitto] · 0.95"},
					{"resa", "11 jun 08:21"},
					{"åtgärd", "→ kalender (förslag)"},
				},
				Tasks: []string{"Lägg i kalender"},
			},
		},
		{
			ID: "m10", Cat: "intressant", From: "Hacker News Daily", FromAddr: "digest@hndaily.com",
			Avatar: "HN", Subject: "Top: \"Show HN: I built a vim-first mail client\"",
			Snippet: "847 poäng · 312 kommentarer. En diskussion om keyboard-driven e-post och varför …",
			Time:    "ons", Day: "tidigare", Unread: false, Confidence: 0.78, ThreadCount: 1,
			Body: []string{"847 poäng, 312 kommentarer. Tråd om keyboard-first mail."},
			Agent: AgentAnalysis{
				Summary: "HN-digest med en tråd som matchar ditt projekt. Topplänk relevant.",
				Facts: []Fact{
					{"kategori", "[intressant] · 0.78"},
					{"matchar", "projekt: mailklient"},
				},
				Tasks: []string{"Spara länk"},
			},
		},
		{
			ID: "m11", Cat: "faktura", From: "Telia Företag", FromAddr: "faktura@telia.se",
			Avatar: "Te", Subject: "Månadsfaktura juni",
			Snippet: "Belopp 689 kr. Autogiro dras 27 jun. Ingen åtgärd krävs …",
			Time:    "tis", Day: "tidigare", Unread: false, Confidence: 0.96, ThreadCount: 1,
			Body: []string{"Belopp: 689 kr\nDras via autogiro 27 jun. Ingen åtgärd krävs."},
			Agent: AgentAnalysis{
				Summary: "Faktura 689 kr, autogiro — ingen åtgärd. Bokförs automatiskt.",
				Facts: []Fact{
					{"kategori", "[faktura] · 0.96"},
					{"belopp", "689 kr · autogiro"},
					{"åtgärd", "ingen — informativ"},
				},
				Tasks: []string{},
			},
		},
		{
			ID: "m12", Cat: "reklam", From: "Linear", FromAddr: "hello@linear.app",
			Avatar: "LN", Subject: "Changelog: snabbare sökning + nya kortkommandon",
			Snippet: "Den här veckan: 3× snabbare global sökning, fler kortkommandon och …",
			Time:    "tis", Day: "tidigare", Unread: false, Confidence: 0.90, ThreadCount: 1,
			Body: []string{"Veckans changelog från Linear."},
			Agent: AgentAnalysis{
				Summary: "Produktchangelog. Auto-arkiveras enligt regel.",
				Facts: []Fact{
					{"kategori", "[reklam] · 0.90"},
					{"regel", "auto-arkivera changelog"},
				},
				Tasks: []string{},
			},
		},
	}
}

func mockActivity() []Activity {
	return []Activity{
		{Time: "09:13", Text: "klassificerade 4 nya mail"},
		{Time: "09:13", Text: "genererade utkast → Anna Lindqvist", Cat: "svara"},
		{Time: "09:12", Text: "auto-arkiverade Figma (regel: nyhetsbrev)", Cat: "reklam"},
		{Time: "08:41", Text: "flaggade säkerhetsärende → David Ek", Cat: "svara"},
		{Time: "07:55", Text: "lade Fortnox-faktura i betalkö", Cat: "faktura"},
	}
}

func mockRules() []Rule {
	return []Rule{
		{On: true, Text: "Auto-arkivera nyhetsbrev från kända avsändare", Scope: "reklam"},
		{On: true, Text: "Drafta svar automatiskt för [att svara på]", Scope: "svara"},
		{On: true, Text: "Lägg fakturor i betalkö och matcha mot underlag", Scope: "faktura"},
		{On: true, Text: "Skicka aldrig utan mitt godkännande", Scope: "global", Locked: true},
		{On: true, Text: "Föreslå kalenderhändelser från biljetter & bokningar", Scope: "kvitto"},
		{On: false, Text: "Sammanfatta trådar längre än 6 mail", Scope: "global"},
	}
}

func mockAccounts() []Account {
	return []Account{
		{
			Email:    "johan@nordveda.se",
			Provider: "Gmail",
			Status:   StatusConnected,
			SyncedAt: "synkad nyss",
			Scopes: []string{
				"läser e-post",
				"skapar utkast, skickar aldrig automatiskt",
				"hanterar etiketter",
			},
		},
	}
}
