package main

func panicErr(err error){
    if err != nil{
        panic(err)
    }
}
func PanicString(msg string){
    panic(msg)
}