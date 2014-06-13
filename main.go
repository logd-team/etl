package etl

import (
    "flag"
    "log"
    "os"
    "runtime"
)

func main(){
    runtime.GOMAXPROCS(runtime.NumCPU())

    var (
        outDir string
        spiderList string
        colsMap string
        inChanSize int
        outChanSize int
        saveRoutineNum int
        parseRoutineNum int
        outFilePrefix string
    )
    flag.StringVar(&outDir, "od", "./", "dir to save result")
    flag.StringVar(&spiderList, "sl", "", "spider list file")
    flag.StringVar(&colsMap, "cm", "", "columns map")
    flag.IntVar(&saveRoutineNum, "nrs", 10, "number of routines to save log")
    flag.IntVar(&parseRoutineNum, "nrp", 12, "number of routines to parse log")
    flag.IntVar(&inChanSize, "ics", 100000, "number of input channel size")
    flag.IntVar(&outChanSize, "ocs", 100000, "number of output channel size")
    flag.StringVar(&outFilePrefix, "ofp", "etl", "the output file prefix")
    flag.Parse()

    if spiderList == "" {
        log.Println("[warning] miss spider list")
    }
    if colsMap == "" {
        log.Println("[error] miss columns map file")
        flag.PrintDefaults()
        return
    }
    d := NewDispatcher(colsMap, outDir, saveRoutineNum, outFilePrefix)
    g := NewGlobalHao123(spiderList, inChanSize, outChanSize, parseRoutineNum, d)

    go g.Start(false)

    fnames := flag.Args()
    for _, f := range fnames {
        log.Println("process", f)
        g.ParseFile(f)
    }

    if len(fnames) == 0 {
        g.ParseReader(os.Stdin)
    }

    g.Wait()
}
