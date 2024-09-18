package hw03frequencyanalysis

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Change to true if needed.
var taskWithAsteriskIsCompleted = false

var text = `Как видите, он  спускается  по  лестнице  вслед  за  своим
	другом   Кристофером   Робином,   головой   вниз,  пересчитывая
	ступеньки собственным затылком:  бум-бум-бум.  Другого  способа
	сходить  с  лестницы  он  пока  не  знает.  Иногда ему, правда,
		кажется, что можно бы найти какой-то другой способ, если бы  он
	только   мог   на  минутку  перестать  бумкать  и  как  следует
	сосредоточиться. Но увы - сосредоточиться-то ему и некогда.
		Как бы то ни было, вот он уже спустился  и  готов  с  вами
	познакомиться.
	- Винни-Пух. Очень приятно!
		Вас,  вероятно,  удивляет, почему его так странно зовут, а
	если вы знаете английский, то вы удивитесь еще больше.
		Это необыкновенное имя подарил ему Кристофер  Робин.  Надо
	вам  сказать,  что  когда-то Кристофер Робин был знаком с одним
	лебедем на пруду, которого он звал Пухом. Для лебедя  это  было
	очень   подходящее  имя,  потому  что  если  ты  зовешь  лебедя
	громко: "Пу-ух! Пу-ух!"- а он  не  откликается,  то  ты  всегда
	можешь  сделать вид, что ты просто понарошку стрелял; а если ты
	звал его тихо, то все подумают, что ты  просто  подул  себе  на
	нос.  Лебедь  потом  куда-то делся, а имя осталось, и Кристофер
	Робин решил отдать его своему медвежонку, чтобы оно не  пропало
	зря.
		А  Винни - так звали самую лучшую, самую добрую медведицу
	в  зоологическом  саду,  которую  очень-очень  любил  Кристофер
	Робин.  А  она  очень-очень  любила  его. Ее ли назвали Винни в
	честь Пуха, или Пуха назвали в ее честь - теперь уже никто  не
	знает,  даже папа Кристофера Робина. Когда-то он знал, а теперь
	забыл.
		Словом, теперь мишку зовут Винни-Пух, и вы знаете почему.
		Иногда Винни-Пух любит вечерком во что-нибудь поиграть,  а
	иногда,  особенно  когда  папа  дома,  он больше любит тихонько
	посидеть у огня и послушать какую-нибудь интересную сказку.
		В этот вечер...`

func TestTop10(t *testing.T) {
	t.Run("no words in empty string", func(t *testing.T) {
		require.Len(t, Top10(""), 0)
	})

	t.Run("positive test", func(t *testing.T) {
		if taskWithAsteriskIsCompleted {
			expected := []string{
				"а",         // 8
				"он",        // 8
				"и",         // 6
				"ты",        // 5
				"что",       // 5
				"в",         // 4
				"его",       // 4
				"если",      // 4
				"кристофер", // 4
				"не",        // 4
			}
			require.Equal(t, expected, Top10(text))
		} else {
			expected := []string{
				"он",        // 8
				"а",         // 6
				"и",         // 6
				"ты",        // 5
				"что",       // 5
				"-",         // 4
				"Кристофер", // 4
				"если",      // 4
				"не",        // 4
				"то",        // 4
			}
			require.Equal(t, expected, Top10(text))
		}
	})
}

func TestTop10Hard(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
		hard     bool
	}{
		{
			"Simple words test",
			"Hello, world! This is an example.",
			[]string{
				"Hello,",   // 1
				"This",     // 1
				"an",       // 1
				"example.", // 1
				"is",       // 1
				"world!",   // 1
			},
			false,
		},
		{
			"Simple words hard RE test",
			"Hello, world! This is an example.",
			[]string{
				"an",      // 1
				"example", // 1
				"hello",   // 1
				"is",      // 1
				"this",    // 1
				"world",   // 1
			},
			true,
		},
		{
			"Cyrillic words test",
			"Всем привет! Это пример текста.",
			[]string{
				"Всем",    // 1
				"Это",     // 1
				"привет!", // 1
				"пример",  // 1
				"текста.", // 1
			},
			false,
		},
		{
			"Cyrillic words hard RE test",
			"Всем привет! Это пример текста.",
			[]string{
				"всем",   // 1
				"привет", // 1
				"пример", // 1
				"текста", // 1
				"это",    // 1
			},
			true,
		},
		{
			"Hard phrase test",
			"Hello,world! It's an example. It's a... complex wo...rd 'example'.! Multy-hyphen '----' - is a word too.",
			[]string{
				"It's",         // 2
				"'----'",       // 1
				"'example'.!",  // 1
				"-",            // 1
				"Hello,world!", // 1
				"Multy-hyphen", // 1
				"a",            // 1
				"a...",         // 1
				"an",           // 1
				"complex",      // 1
			},
			false,
		},
		{
			"Hard phrase hard RE test",
			`Hello,world! It's an example. It's a... complex wo...rd 'example'.! Multy-hyphen '----' - is a word too.`,
			[]string{
				"a",            // 2
				"example",      // 2
				"it's",         // 2
				"----",         // 1
				"an",           // 1
				"complex",      // 1
				"hello,world",  // 1
				"is",           // 1
				"multy-hyphen", // 1
				"too",          // 1
			},
			true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := Top10Hard(tc.input, tc.hard)
			require.Equal(t, tc.expected, result)
		})
	}
}
