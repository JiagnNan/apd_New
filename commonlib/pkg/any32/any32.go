package any32

//10进制与32进制互转

var tenToAny map[int64]string = map[int64]string{
	0:  "n",
	1:  "1",
	2:  "m",
	3:  "3",
	4:  "c",
	5:  "v",
	6:  "z",
	7:  "9",
	8:  "a",
	9:  "x",
	10: "t",
	11: "b",
	12: "4",
	13: "d",
	14: "g",
	15: "q",
	16: "e",
	17: "h",
	18: "j",
	19: "k",
	20: "2",
	21: "0",
	22: "p",
	23: "f",
	24: "7",
	25: "s",
	26: "8",
	27: "5",
	28: "w",
	29: "r",
	30: "y",
	31: "6",
}

// 10进制转32进制  不包含 ILOU 字符
func DecimalToAny32(num int64) string {
	n := int64(32)
	str := ""
	for num != 0 {
		remainder := num % n
		remainderStr := tenToAny[remainder]
		str += remainderStr
		num /= n
	}
	return str
}
