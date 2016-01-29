package filemm
import (
    "os"
    "log"
    "bufio"
    "bytes"
    "io"
    "time"
    "math/rand"
)
const (
    NL = "\n"
)

type Job struct {
    line string
    results chan<- Result
    lineNumber int
}
func (this *Job) Do(sign string) {
    count := countSigns(this.line, sign)
    res := Result{line: this.line, sign: sign, count: count, lineNumber: this.lineNumber}
    this.results <- res
}

type Result struct {
    line string
    sign string
    count int
    lineNumber int
}


var workers int = 10

func RunConcurency() {
    rand.Seed(time.Now().Unix())
    args := os.Args
    if len(args) < 2 {
        log.Fatalf("File name not exist. Enter the file name")
        os.Exit(1)
    }
    filename := args[1]
    log.Printf("You have entered a file named %s", filename)
    lines := read(filename)

    jobs := make(chan Job, workers)
    results := make(chan Result, 1000)
    done := make(chan struct{}, workers)
    go addJob(jobs, lines, results)
    for i := 0; i < workers; i++ {
        go doJob(done, "a", jobs)
    }
    awaitCompletion(done, results)
    printResults(results)
}

func awaitCompletion(done <-chan struct{}, results chan<-Result) {
    for i := 0; i < workers; i++ {
        <-done
    }
    close(results)
}

func addJob(jobs chan<- Job, lines []string, results chan<- Result) {
    for pos,line := range lines {
        jobs <- Job{ line, results, pos }
    }
    close(jobs)
}
func doJob(done chan<- struct{}, sign string, jobs <-chan Job) {
    for job := range jobs {
        sleep()
        job.Do(sign)
    }
    done <- struct{}{}
}
func printResults(results <-chan Result) {
    for r := range results {
        log.Printf("%d Result: %s %d", r.lineNumber, r.sign, r.count)
    }
}

func Run() {
    args := os.Args
    if len(args) < 2 {
        log.Fatalf("File name not exist. Enter the file name")
        os.Exit(1)
    }
    filename := args[1]
    log.Printf("You have entered a file named %s", filename)
    lines := read(filename)
    for i := 0; i < len(lines); i++ {
        count := countSigns(lines[i], "a")
        log.Printf("Line %d Count: %d", i, count)
    }
}

func read(filename string) ([]string) {
    file, err := os.Open(filename)
    if err != nil {
        log.Fatalf("Attempt to open a file is unsuccesfull %v", err)
    }
    defer file.Close()
    reader := bufio.NewReader(file)
    slines := make([]string, 0, 10000)
    blines := make([][]byte, 0, 10000)
    for lino := 1; ; lino++ {
        line, err := reader.ReadBytes('\n')
        line = bytes.TrimRight(line, "\n\r")
        if err != nil {
            if err != io.EOF {
                log.Printf("error:%d: %s\n", lino, err)
            }
            break
        }
        slines = append(slines, string(line))
        blines = append(blines, line)
    }
    return slines
}

func sleep() {
    randomNumber := rand.Intn(1000 - 500) + 500
    log.Printf("count signs %d", randomNumber)
    time.Sleep(time.Duration(randomNumber) * time.Millisecond)
}

func countSigns(line, sign string) int {
    var count int = 0
    for i := 0; i < len(line); i++ {
        if line[i] == sign[0]    {
            count++
        }
    }
    return count
}
