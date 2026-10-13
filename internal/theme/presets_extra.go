package theme

// extraPresets rounds the library out to cover every venue type. They are
// appended to the base set in Presets().
//
// Each builder starts from Default() and overrides only what it needs, so a new
// field added to Config automatically gets a sane value in every preset.
var extraPresets = []Preset{
	// ---------------------------------------------------------------- cafés
	{
		Key: "latte", LabelFa: "لاته", LabelEn: "Latte",
		NoteFa: "کرم روشن و آرام، مناسب کافه", NoteEn: "Soft cream tones for cafés",
		For:    []string{"cafe", "bakery", "juicebar"},
		Swatch: []string{"#f7f2ea", "#8c6a4f", "#3f3630"},
		Build: func() Config {
			c := Default()
			c.Colors = Colors{
				Primary: "#8c6a4f", Secondary: "#6b5340", Accent: "#c39b73",
				Background: "#f7f2ea", Surface: "#fffdf9", SurfaceAlt: "#f0e7da",
				Text: "#3f3630", TextMuted: "#8a7a6c", Border: "#e3d7c6",
				Price: "#8c6a4f", Badge: "#c39b73",
			}
			c.Background.Type = "solid"
			c.Background.Color = "#f7f2ea"
			c.Typography.FontEn = "playfair"
			c.Typography.LineHeight = 1.85
			c.Cards.Style = "flat"
			c.Cards.Radius = 22
			c.Cards.Shadow = 1
			c.Header.Style = "centered"
			c.Layout.Density = "relaxed"
			return c
		},
	},
	{
		Key: "espresso", LabelFa: "اسپرسو", LabelEn: "Espresso",
		NoteFa: "تیره و فشرده، کافه مدرن", NoteEn: "Dark and tight, modern coffee bar",
		For:    []string{"cafe", "bakery"},
		Swatch: []string{"#191512", "#d9a441", "#efe6da"},
		Build: func() Config {
			c := Default()
			c.Colors = Colors{
				Primary: "#d9a441", Secondary: "#2c241d", Accent: "#b5764a",
				Background: "#191512", Surface: "#221c17", SurfaceAlt: "#2b231c",
				Text: "#efe6da", TextMuted: "#a2917f", Border: "#352b23",
				Price: "#d9a441", Badge: "#b5764a",
			}
			c.Background.Type = "solid"
			c.Background.Color = "#191512"
			c.Cards.Style = "outline"
			c.Cards.Radius = 6
			c.Layout.Mode = "compact"
			c.Layout.CategoryNav = "tabs"
			c.Typography.LetterSpacing = 0.3
			c.Effects.RoundedBtns = false
			return c
		},
	},
	{
		Key: "matcha", LabelFa: "ماچا", LabelEn: "Matcha",
		NoteFa: "سبز ملایم و تازه", NoteEn: "Gentle fresh green",
		For:    []string{"cafe", "juicebar", "icecream"},
		Swatch: []string{"#f2f7f0", "#4f7942", "#d98555"},
		Build: func() Config {
			c := Default()
			c.Colors = Colors{
				Primary: "#4f7942", Secondary: "#37592f", Accent: "#d98555",
				Background: "#f2f7f0", Surface: "#ffffff", SurfaceAlt: "#e6efe2",
				Text: "#26301f", TextMuted: "#69775f", Border: "#d5e3ce",
				Price: "#4f7942", Badge: "#d98555",
			}
			c.Background.Type = "gradient"
			c.Background.GradientFrom = "#f6faf4"
			c.Background.GradientVia = "#eef5eb"
			c.Background.GradientTo = "#e3eede"
			c.Background.GradientDeg = 175
			c.Cards.Radius = 26
			c.Cards.Shadow = 1
			c.Layout.Columns = 3
			return c
		},
	},

	// ------------------------------------------------------------ fast food
	{
		Key: "burger", LabelFa: "برگر", LabelEn: "Burger",
		NoteFa: "قرمز و زرد پرانرژی", NoteEn: "Bold red and yellow",
		For:    []string{"fastfood"},
		Swatch: []string{"#1b1b1f", "#e63946", "#ffbe0b"},
		Build: func() Config {
			c := Default()
			c.Colors = Colors{
				Primary: "#e63946", Secondary: "#2a2a30", Accent: "#ffbe0b",
				Background: "#1b1b1f", Surface: "#25252b", SurfaceAlt: "#2f2f36",
				Text: "#f7f7f8", TextMuted: "#a3a3ad", Border: "#3a3a43",
				Price: "#ffbe0b", Badge: "#e63946",
			}
			c.Background.Type = "solid"
			c.Background.Color = "#1b1b1f"
			c.Typography.HeadingWeight = 900
			c.Typography.FontEn = "montserrat"
			c.Cards.Style = "elevated"
			c.Cards.Radius = 18
			c.Cards.Shadow = 3
			c.Cards.HoverEffect = "zoom"
			c.Layout.Columns = 3
			return c
		},
	},
	{
		Key: "pizzeria", LabelFa: "پیتزریا", LabelEn: "Pizzeria",
		NoteFa: "قرمز و سبز ایتالیایی", NoteEn: "Italian red and green",
		For:    []string{"fastfood", "restaurant"},
		Swatch: []string{"#fbf7f0", "#c1272d", "#2e7d32"},
		Build: func() Config {
			c := Default()
			c.Colors = Colors{
				Primary: "#c1272d", Secondary: "#2e7d32", Accent: "#e8a33d",
				Background: "#fbf7f0", Surface: "#ffffff", SurfaceAlt: "#f3ece0",
				Text: "#2b2320", TextMuted: "#7d7169", Border: "#e6dbc9",
				Price: "#c1272d", Badge: "#2e7d32",
			}
			c.Background.Type = "pattern"
			c.Background.Color = "#fbf7f0"
			c.Background.Pattern = "crosshatch"
			c.Background.PatternOpac = 0.45
			c.Background.PatternScale = 18
			c.Typography.FontEn = "playfair"
			c.Cards.Radius = 12
			c.Cards.Shadow = 2
			return c
		},
	},
	{
		Key: "arcade", LabelFa: "آرکید", LabelEn: "Arcade",
		NoteFa: "رترو و پرزرق‌وبرق", NoteEn: "Retro and loud",
		For:    []string{"fastfood", "icecream"},
		Swatch: []string{"#170b2e", "#ff3cac", "#00f5d4"},
		Build: func() Config {
			c := Default()
			c.Colors = Colors{
				Primary: "#ff3cac", Secondary: "#2b1a55", Accent: "#00f5d4",
				Background: "#170b2e", Surface: "#221142", SurfaceAlt: "#2c1754",
				Text: "#f6efff", TextMuted: "#a894cc", Border: "#3d2270",
				Price: "#00f5d4", Badge: "#ff3cac",
			}
			c.Background.Type = "gradient"
			c.Background.GradientFrom = "#12081f"
			c.Background.GradientVia = "#26104a"
			c.Background.GradientTo = "#4a1a63"
			c.Background.GradientDeg = 145
			c.Cards.Style = "glass"
			c.Cards.HoverEffect = "glow"
			c.Cards.Radius = 20
			c.Effects.Glass = true
			c.Typography.HeadingWeight = 900
			return c
		},
	},

	// -------------------------------------------------- Persian & grill houses
	{
		Key: "ember", LabelFa: "اخگر", LabelEn: "Ember",
		NoteFa: "زغالی و آتشین، مناسب جگرکی", NoteEn: "Charcoal and fire, for grill houses",
		For:    []string{"jigaraki", "restaurant", "fastfood"},
		Swatch: []string{"#120e0c", "#ff6b35", "#f7c59f"},
		Build: func() Config {
			c := Default()
			c.Colors = Colors{
				Primary: "#ff6b35", Secondary: "#2a1d17", Accent: "#f7c59f",
				Background: "#120e0c", Surface: "#1c1613", SurfaceAlt: "#251c17",
				Text: "#f6ece5", TextMuted: "#a8968a", Border: "#332721",
				Price: "#ff6b35", Badge: "#f7c59f",
			}
			c.Background.Type = "gradient"
			c.Background.GradientFrom = "#0b0807"
			c.Background.GradientVia = "#160f0c"
			c.Background.GradientTo = "#2b1610"
			c.Background.GradientDeg = 190
			c.Cards.Radius = 10
			c.Cards.Shadow = 3
			c.Cards.HoverEffect = "glow"
			c.Typography.HeadingWeight = 800
			return c
		},
	},
	{
		Key: "turquoise", LabelFa: "فیروزه", LabelEn: "Turquoise",
		NoteFa: "فیروزه‌ای اصفهانی با طلایی", NoteEn: "Persian turquoise and gold",
		For:    []string{"restaurant", "teahouse", "cafe"},
		Swatch: []string{"#07333a", "#3fbfb5", "#e3b23c"},
		Build: func() Config {
			c := Default()
			c.Colors = Colors{
				Primary: "#3fbfb5", Secondary: "#0c4952", Accent: "#e3b23c",
				Background: "#07333a", Surface: "#0c4149", SurfaceAlt: "#104d56",
				Text: "#eafaf8", TextMuted: "#8fb9b8", Border: "#16606b",
				Price: "#e3b23c", Badge: "#3fbfb5",
			}
			c.Background.Type = "pattern"
			c.Background.Color = "#07333a"
			c.Background.Pattern = "diagonal"
			c.Background.PatternOpac = 0.5
			c.Background.PatternScale = 30
			c.Typography.FontFa = "nazanin"
			c.Cards.Radius = 16
			c.Header.Style = "centered"
			return c
		},
	},
	{
		Key: "clay", LabelFa: "سفالی", LabelEn: "Clay",
		NoteFa: "خاکی و سنتی، سفره‌خانه", NoteEn: "Earthy and traditional, teahouse",
		For:    []string{"teahouse", "restaurant", "jigaraki"},
		Swatch: []string{"#f4ece1", "#a0522d", "#4a5240"},
		Build: func() Config {
			c := Default()
			c.Colors = Colors{
				Primary: "#a0522d", Secondary: "#4a5240", Accent: "#c98a5e",
				Background: "#f4ece1", Surface: "#fdf9f3", SurfaceAlt: "#eae0d1",
				Text: "#33291f", TextMuted: "#7f7263", Border: "#ddd0bd",
				Price: "#a0522d", Badge: "#4a5240",
			}
			c.Background.Type = "pattern"
			c.Background.Color = "#f4ece1"
			c.Background.Pattern = "dots"
			c.Background.PatternOpac = 0.6
			c.Background.PatternScale = 16
			c.Typography.FontFa = "nazanin"
			c.Typography.HeadingWeight = 700
			c.Cards.Style = "outline"
			c.Cards.Radius = 8
			c.Layout.Mode = "list"
			return c
		},
	},
	{
		Key: "nightgarden", LabelFa: "باغ شب", LabelEn: "Night Garden",
		NoteFa: "سبز تیره با طلایی، قهوه‌خانه", NoteEn: "Deep green and gold, teahouse",
		For:    []string{"teahouse", "restaurant", "cafe"},
		Swatch: []string{"#101b16", "#c9a227", "#7d9d6f"},
		Build: func() Config {
			c := Default()
			c.Colors = Colors{
				Primary: "#c9a227", Secondary: "#1c2e25", Accent: "#7d9d6f",
				Background: "#101b16", Surface: "#18271f", SurfaceAlt: "#1f3228",
				Text: "#eef4ee", TextMuted: "#9aae9d", Border: "#284034",
				Price: "#c9a227", Badge: "#7d9d6f",
			}
			c.Background.Type = "gradient"
			c.Background.GradientFrom = "#0a120e"
			c.Background.GradientVia = "#101b16"
			c.Background.GradientTo = "#17281f"
			c.Background.GradientDeg = 200
			c.Typography.FontEn = "playfair"
			c.Cards.Radius = 18
			c.Cards.Style = "glass"
			c.Effects.Glass = true
			return c
		},
	},

	// --------------------------------------------------- bakery & ice cream
	{
		Key: "vanilla", LabelFa: "وانیل", LabelEn: "Vanilla",
		NoteFa: "کرم و صورتی ملایم، شیرینی", NoteEn: "Cream and soft pink, pastry",
		For:    []string{"bakery", "icecream", "cafe"},
		Swatch: []string{"#fff9f2", "#d98da6", "#8a6f4e"},
		Build: func() Config {
			c := Default()
			c.Colors = Colors{
				Primary: "#d98da6", Secondary: "#8a6f4e", Accent: "#f0b6a0",
				Background: "#fff9f2", Surface: "#ffffff", SurfaceAlt: "#fdeee6",
				Text: "#4a3b34", TextMuted: "#94806f", Border: "#f2e0d4",
				Price: "#d98da6", Badge: "#f0b6a0",
			}
			c.Background.Type = "solid"
			c.Background.Color = "#fff9f2"
			c.Typography.FontEn = "poppins"
			c.Cards.Radius = 30
			c.Cards.Shadow = 1
			c.Cards.ImageRatio = "square"
			c.Layout.Columns = 3
			c.Layout.Density = "relaxed"
			return c
		},
	},
	{
		Key: "sorbet", LabelFa: "شربتی", LabelEn: "Sorbet",
		NoteFa: "رنگ‌های شاد پاستلی، بستنی", NoteEn: "Bright pastels, ice cream",
		For:    []string{"icecream", "juicebar", "cafe"},
		Swatch: []string{"#fef6fb", "#ff8fab", "#4cc9f0"},
		Build: func() Config {
			c := Default()
			c.Colors = Colors{
				Primary: "#ff8fab", Secondary: "#4cc9f0", Accent: "#ffd166",
				Background: "#fef6fb", Surface: "#ffffff", SurfaceAlt: "#fdeaf3",
				Text: "#40323a", TextMuted: "#93818b", Border: "#f7dce8",
				Price: "#ff8fab", Badge: "#4cc9f0",
			}
			c.Background.Type = "gradient"
			c.Background.GradientFrom = "#fef6fb"
			c.Background.GradientVia = "#f4f7fe"
			c.Background.GradientTo = "#eefaf8"
			c.Background.GradientDeg = 160
			c.Cards.Radius = 32
			c.Cards.Shadow = 2
			c.Cards.HoverEffect = "zoom"
			c.Cards.ImageRatio = "square"
			c.Layout.Columns = 3
			c.Typography.FontEn = "poppins"
			return c
		},
	},
	{
		Key: "citrus", LabelFa: "مرکبات", LabelEn: "Citrus",
		NoteFa: "نارنجی و سبز تازه، آبمیوه", NoteEn: "Fresh orange and green, juice bar",
		For:    []string{"juicebar", "icecream", "fastfood"},
		Swatch: []string{"#fffdf0", "#f77f00", "#43aa8b"},
		Build: func() Config {
			c := Default()
			c.Colors = Colors{
				Primary: "#f77f00", Secondary: "#43aa8b", Accent: "#fcbf49",
				Background: "#fffdf0", Surface: "#ffffff", SurfaceAlt: "#fdf6dc",
				Text: "#2f2a1c", TextMuted: "#83795f", Border: "#efe6c4",
				Price: "#f77f00", Badge: "#43aa8b",
			}
			c.Background.Type = "solid"
			c.Background.Color = "#fffdf0"
			c.Cards.Radius = 24
			c.Cards.Shadow = 2
			c.Cards.ImageRatio = "square"
			c.Layout.Columns = 3
			c.Typography.HeadingWeight = 800
			return c
		},
	},

	// -------------------------------------------------------- universal looks
	{
		Key: "paper", LabelFa: "کاغذی", LabelEn: "Paper",
		NoteFa: "شبیه منوی چاپی، بدون عکس", NoteEn: "Like a printed menu, no photos",
		Swatch: []string{"#faf8f3", "#1f1c17", "#8b7355"},
		Build: func() Config {
			c := Default()
			c.Colors = Colors{
				Primary: "#1f1c17", Secondary: "#8b7355", Accent: "#8b7355",
				Background: "#faf8f3", Surface: "#faf8f3", SurfaceAlt: "#f2eee4",
				Text: "#1f1c17", TextMuted: "#7a705f", Border: "#ddd6c6",
				Price: "#1f1c17", Badge: "#8b7355",
			}
			c.Background.Type = "solid"
			c.Background.Color = "#faf8f3"
			c.Typography.FontEn = "georgia"
			c.Typography.FontFa = "nazanin"
			c.Typography.HeadingScale = 1.32
			c.Cards.Style = "flat"
			c.Cards.Radius = 0
			c.Cards.Shadow = 0
			c.Cards.BorderWidth = 0
			c.Cards.ImagePos = "none"
			c.Cards.HoverEffect = "none"
			c.Layout.Mode = "list"
			c.Layout.CategoryNav = "none"
			c.Layout.MaxWidth = 760
			c.Display.ShowImages = false
			c.Header.Style = "centered"
			c.Effects.Animations = false
			c.Effects.RoundedBtns = false
			return c
		},
	},
	{
		Key: "brutalist", LabelFa: "بروتالیست", LabelEn: "Brutalist",
		NoteFa: "خطوط ضخیم و کنتراست بالا", NoteEn: "Heavy borders, high contrast",
		Swatch: []string{"#ffffff", "#000000", "#ff4d00"},
		Build: func() Config {
			c := Default()
			c.Colors = Colors{
				Primary: "#000000", Secondary: "#ff4d00", Accent: "#ff4d00",
				Background: "#ffffff", Surface: "#ffffff", SurfaceAlt: "#f2f2f2",
				Text: "#000000", TextMuted: "#555555", Border: "#000000",
				Price: "#ff4d00", Badge: "#000000",
			}
			c.Background.Type = "solid"
			c.Background.Color = "#ffffff"
			c.Cards.Style = "outline"
			c.Cards.Radius = 0
			c.Cards.Shadow = 0
			c.Cards.BorderWidth = 3
			c.Cards.HoverEffect = "lift"
			c.Typography.HeadingWeight = 900
			c.Typography.FontEn = "montserrat"
			c.Typography.LetterSpacing = -0.3
			c.Layout.CategoryNav = "tabs"
			c.Effects.RoundedBtns = false
			return c
		},
	},
	{
		Key: "aurora", LabelFa: "شفق", LabelEn: "Aurora",
		NoteFa: "گرادیان رنگی و شیشه‌ای", NoteEn: "Colourful gradient with glass",
		Swatch: []string{"#0d1030", "#7f5af0", "#2cb67d"},
		Build: func() Config {
			c := Default()
			c.Colors = Colors{
				Primary: "#7f5af0", Secondary: "#232946", Accent: "#2cb67d",
				Background: "#0d1030", Surface: "#171a3a", SurfaceAlt: "#1f2347",
				Text: "#fffffe", TextMuted: "#94a1b2", Border: "#2b3160",
				Price: "#2cb67d", Badge: "#7f5af0",
			}
			c.Background.Type = "gradient"
			c.Background.GradientFrom = "#0d1030"
			c.Background.GradientVia = "#2b1b62"
			c.Background.GradientTo = "#12403f"
			c.Background.GradientDeg = 140
			c.Cards.Style = "glass"
			c.Cards.Radius = 26
			c.Cards.HoverEffect = "lift"
			c.Effects.Glass = true
			c.Typography.FontEn = "poppins"
			return c
		},
	},
	{
		Key: "monoline", LabelFa: "تک‌خط", LabelEn: "Monoline",
		NoteFa: "خاکستری آرام و بسیار ساده", NoteEn: "Calm grey, very restrained",
		Swatch: []string{"#fafafa", "#3f3f46", "#0ea5e9"},
		Build: func() Config {
			c := Default()
			c.Colors = Colors{
				Primary: "#3f3f46", Secondary: "#71717a", Accent: "#0ea5e9",
				Background: "#fafafa", Surface: "#ffffff", SurfaceAlt: "#f4f4f5",
				Text: "#18181b", TextMuted: "#71717a", Border: "#e4e4e7",
				Price: "#0ea5e9", Badge: "#3f3f46",
			}
			c.Background.Type = "solid"
			c.Background.Color = "#fafafa"
			c.Cards.Style = "flat"
			c.Cards.Radius = 10
			c.Cards.Shadow = 0
			c.Cards.BorderWidth = 1
			c.Cards.ImageRatio = "16:9"
			c.Layout.Mode = "list"
			c.Layout.Density = "compact"
			c.Typography.HeadingWeight = 600
			return c
		},
	},
	{
		Key: "sandstone", LabelFa: "ماسه‌سنگ", LabelEn: "Sandstone",
		NoteFa: "بژ گرم و آفتابی", NoteEn: "Warm sunlit beige",
		Swatch: []string{"#f5efe6", "#b45309", "#57534e"},
		Build: func() Config {
			c := Default()
			c.Colors = Colors{
				Primary: "#b45309", Secondary: "#57534e", Accent: "#d97706",
				Background: "#f5efe6", Surface: "#fffbf5", SurfaceAlt: "#ece2d4",
				Text: "#3b322a", TextMuted: "#857567", Border: "#ded1bf",
				Price: "#b45309", Badge: "#d97706",
			}
			c.Background.Type = "gradient"
			c.Background.GradientFrom = "#f8f3ea"
			c.Background.GradientVia = "#f2ebdf"
			c.Background.GradientTo = "#e9dfcd"
			c.Background.GradientDeg = 170
			c.Cards.Radius = 20
			c.Cards.Shadow = 1
			c.Layout.Mode = "magazine"
			return c
		},
	},
	{
		Key: "obsidian", LabelFa: "اوبسیدین", LabelEn: "Obsidian",
		NoteFa: "سیاه مطلق با نقره‌ای", NoteEn: "Pure black with silver",
		Swatch: []string{"#000000", "#c0c0c0", "#8b5cf6"},
		Build: func() Config {
			c := Default()
			c.Colors = Colors{
				Primary: "#d4d4d8", Secondary: "#18181b", Accent: "#8b5cf6",
				Background: "#000000", Surface: "#0b0b0d", SurfaceAlt: "#131316",
				Text: "#fafafa", TextMuted: "#8b8b93", Border: "#232327",
				Price: "#d4d4d8", Badge: "#8b5cf6",
			}
			c.Background.Type = "solid"
			c.Background.Color = "#000000"
			c.Cards.Style = "outline"
			c.Cards.Radius = 4
			c.Cards.BorderWidth = 1
			c.Cards.HoverEffect = "glow"
			c.Typography.LetterSpacing = 0.6
			c.Typography.HeadingWeight = 500
			c.Effects.RoundedBtns = false
			return c
		},
	},
	{
		Key: "coral", LabelFa: "مرجانی", LabelEn: "Coral",
		NoteFa: "مرجانی و آبی، غذای دریایی", NoteEn: "Coral and blue, seafood",
		For:    []string{"restaurant", "fastfood", "juicebar"},
		Swatch: []string{"#f0fbff", "#ff6f61", "#1d3557"},
		Build: func() Config {
			c := Default()
			c.Colors = Colors{
				Primary: "#ff6f61", Secondary: "#1d3557", Accent: "#457b9d",
				Background: "#f0fbff", Surface: "#ffffff", SurfaceAlt: "#e2f2f8",
				Text: "#12293d", TextMuted: "#5f7a8c", Border: "#cee7f0",
				Price: "#ff6f61", Badge: "#457b9d",
			}
			c.Background.Type = "gradient"
			c.Background.GradientFrom = "#f4fcff"
			c.Background.GradientVia = "#eaf6fb"
			c.Background.GradientTo = "#dcf0f7"
			c.Background.GradientDeg = 185
			c.Cards.Radius = 22
			c.Cards.Shadow = 2
			return c
		},
	},
	{
		Key: "velvet", LabelFa: "مخملی", LabelEn: "Velvet",
		NoteFa: "قرمز مخملی و لوکس", NoteEn: "Rich velvet red",
		Swatch: []string{"#1a0508", "#a4133c", "#ffb3c1"},
		Build: func() Config {
			c := Default()
			c.Colors = Colors{
				Primary: "#a4133c", Secondary: "#590d22", Accent: "#ffb3c1",
				Background: "#1a0508", Surface: "#26080e", SurfaceAlt: "#330b14",
				Text: "#fff0f3", TextMuted: "#c9909d", Border: "#4a1220",
				Price: "#ffb3c1", Badge: "#a4133c",
			}
			c.Background.Type = "gradient"
			c.Background.GradientFrom = "#120406"
			c.Background.GradientVia = "#1f060c"
			c.Background.GradientTo = "#3a0a17"
			c.Background.GradientDeg = 165
			c.Typography.FontEn = "playfair"
			c.Typography.HeadingScale = 1.38
			c.Cards.Style = "gradient"
			c.Cards.Radius = 20
			c.Cards.Shadow = 4
			c.Header.Style = "hero"
			return c
		},
	},
	{
		Key: "linen", LabelFa: "کتان", LabelEn: "Linen",
		NoteFa: "سفید مات و بسیار تمیز", NoteEn: "Matte white, very clean",
		Swatch: []string{"#f6f6f4", "#2f4f4f", "#a67c52"},
		Build: func() Config {
			c := Default()
			c.Colors = Colors{
				Primary: "#2f4f4f", Secondary: "#a67c52", Accent: "#a67c52",
				Background: "#f6f6f4", Surface: "#ffffff", SurfaceAlt: "#eeeeea",
				Text: "#2b2b28", TextMuted: "#7b7b72", Border: "#e0e0d8",
				Price: "#2f4f4f", Badge: "#a67c52",
			}
			c.Background.Type = "solid"
			c.Background.Color = "#f6f6f4"
			c.Cards.Style = "flat"
			c.Cards.Radius = 14
			c.Cards.Shadow = 1
			c.Cards.ImageRatio = "4:3"
			c.Layout.Density = "relaxed"
			c.Layout.MaxWidth = 1040
			c.Typography.LineHeight = 1.9
			return c
		},
	},
}
