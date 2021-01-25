package lexer

type RuneTape struct {
	Channel chan rune
	buffer  []rune
	index   int
}

func NewRuneTape(input chan rune) *RuneTape {
	return &RuneTape{
		Channel: input,
		buffer:  []rune{},
		index:   0,
	}
}

func (t *RuneTape) append(runes ...rune) {
	t.buffer = append(t.buffer, runes...)
}

func (t *RuneTape) peek() rune {
	return t.runeAt(t.index)
}

func (t *RuneTape) runeAt(index int) rune {
	if index >= len(t.buffer) {
		required := index - len(t.buffer) + 1
		t.readFromChannel(required)
	}
	return t.buffer[index]
}

func (t *RuneTape) advance() rune {
	next := t.runeAt(t.index)
	t.index++
	return next
}

func (t *RuneTape) readFromChannel(amount int) {
	for amount > 0 {
		amount--
		nextRune := <-t.Channel
		t.append(nextRune)
	}
}
