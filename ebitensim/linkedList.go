package ebitensim

/*
import "log"

func compareHeight(a, b *ebitenSprite) int {
	if a.z > b.z {
		return 1
	}
	if a.z < b.z {
		return -1
	}

	// A lower y is higher
	if a.y <= b.y {
		return 1
	}
	return -1
}


func llInsert(head **ebitenSprite, new *ebitenSprite) {
	if head == nil {
		log.Println("head is empty")
		return
	}

	if *head == nil {
		*head = new
		return
	}

	// Make the lower one first
	if compareHeight(new, *head) < 0 {
		h := *head
		h.llPrev = new

		*head = new
		new.llNext = h
		new.llPrev = nil
		return
	}

	prev := *head
	next := (*head).llNext
	for next != nil {
		if compareHeight(new, next) < 0 {
			prev.llNext = new
			new.llPrev = prev
			new.llNext = next

			next.llPrev = new
			return
		}

		prev = next
		next = next.llNext
	}

	prev.llNext = new
	new.llPrev = prev
	new.llNext = nil
}

func llRemove(head **ebitenSprite, in *ebitenSprite) {
	if head == nil || *head == nil {
		log.Println("head is empty")
		return
	}

	if *head == in {
		*head = in.llNext
		return
	}

	prev := in.llPrev
	if prev == nil {
		return // Not actually in the LL
	}

	prev.llNext = in.llNext
}
*/