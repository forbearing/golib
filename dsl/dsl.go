package dsl

func Create(func())             {}
func Delete(func())             {}
func Update(func())             {}
func UpdatePartial(func())      {}
func Get(func())                {}
func List(func())               {}
func BatchCreate(func())        {}
func BatchDelete(func())        {}
func BatchUpdate(func())        {}
func BatchUpdatePartial(func()) {}

func Payload(any) {}
func Result(any)  {}

func Endpoint(string) {}
func Enabled(bool)    {}
