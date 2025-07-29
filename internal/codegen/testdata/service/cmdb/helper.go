package service_cmdb

func fillRow(row *[]string, expectLen int) {
	for i := len(*row); i < expectLen; i++ {
		*row = append(*row, "")
	}
}
