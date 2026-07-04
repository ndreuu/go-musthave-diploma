package service

func IsValidLuhn(number string) bool {
	if number == "" {
		return false
	}

	sum := 0
	alternate := false

	for i := len(number) - 1; i >= 0; i-- {
		if number[i] < '0' || number[i] > '9' {
			return false
		}

		n := int(number[i] - '0')

		if alternate {
			n *= 2
			if n > 9 {
				n -= 9
			}
		}

		sum += n
		alternate = !alternate
	}

	return sum%10 == 0
}
