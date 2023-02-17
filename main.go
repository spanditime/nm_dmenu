package main

import (
    "strconv"
    "sort"
    "os/exec"
    "strings"
)

type connection struct{
    name string
    uuid string
    ctype string
    device string
    active bool
}

type wifi struct{
    active bool
    BSSID string
    SSID string
    mode string
    ch string
    rate string
    signal string
    bars string
    security string
}

var savedConections []connection
var savedWifis []wifi

func nmcli(args []string, wait bool) (string, error) {
    pargs := []string{"-t"}
    if !wait {
        pargs = append(pargs, []string{"-w", "0"}...)
    }
    args = append(pargs, args...)
    out, err := exec.Command("nmcli", args...).Output()
    if err != nil {
        return "", err
    }
    return string(out[:]), err
}

func turnNetworkingOn(wait bool) error {
    _, err := nmcli([]string{"n","on"},wait)
    return err
}

func turnNetworkingOff(wait bool) error {
    _, err := nmcli([]string{"n","off"},wait)
    return err
}

func toggleNetworking(wait bool) (error) {
    netw, err := getNetworking()
    if err != nil{
        return err
    }
    if netw {
        err = turnNetworkingOff(wait)
    }else{
        err = turnNetworkingOn(wait)
    }
    return err
}

func getNetworking() (bool, error){
    networking := false
    out, err := nmcli([]string{"n"},false)
    if err != nil {
        return networking, nil
    }
    if out == "enabled\n"{
        networking = true
    }
    return networking, nil
}

func connectBySSIDorUUID(SSIDorUUID string) error {
    _,err := nmcli([]string{"c","u", SSIDorUUID},false)
    return err
}

func disconnectBySSIDorUUID(SSIDorUUID string) error {
    _,err := nmcli([]string{"c","d", SSIDorUUID},false)
    return err
}

func connectWiFiBySSIDorBSSID(SSIDorBSSID string) error {
    _, err := nmcli([]string{"d","w","c",SSIDorBSSID},false)
    return err
}


func rescanWiFis(wait bool) error {
    _, err := nmcli([]string{"d","w","r"},wait)
    return err
}

func stringToCons(str string) []connection {
    cons := []connection{}
    constrs := strings.Split(str,"\n")
    for _, constr := range constrs {
        con := connection{}
        conparams := strings.Split(constr, ":")
        if len(conparams) < 3 {
            continue
        }
        con.name = conparams[0]
        con.uuid = conparams[1]
        con.ctype = conparams[2]
        if len(conparams) >= 4 {
            con.device = conparams[3]
        }
        cons = append(cons, con)
    }
    return cons
}

func getSavedCons() ([]connection,error) {
    out,err := nmcli([]string{"c","s"},false)
    if err != nil {
        return nil, err
    }
    savedcons := stringToCons(out)
    out,err = nmcli([]string{"c","s","--active"},false)
    if err != nil {
        return nil, err
    }
    activecons := stringToCons(out)
    for j:=0; j<len(savedcons); j++{
        if len(activecons) == 0 {
            break;
        }
        for i := 0; i < len(activecons); i++ {
            if savedcons[j].uuid != activecons[i].uuid{
                continue
            }
            savedcons[j].active = true
            activecons = append(activecons[:i], activecons[i+1:]...)
            i--
        }
    }
    return savedcons,nil
}

func stringToWiFis(str string) []wifi {
    wifis := []wifi{}
    wifistrs := strings.Split(str,"\n")
    for _, wifistr := range wifistrs{
        wifiparams := strings.Split(wifistr,":")
        wifi := wifi{}
        if len(wifiparams) < 13 {
            continue
        }
        wifi.active = wifiparams[0] == "*"
        for i:=1; i<6; i++{
            wifi.BSSID += strings.Replace(wifiparams[i],"\\","",1) + ":"
        }
        wifi.BSSID += wifiparams[6]
        wifi.SSID = wifiparams[7]
        wifi.mode = wifiparams[8]
        wifi.ch = wifiparams[9]
        wifi.rate = wifiparams[10]
        wifi.signal = wifiparams[11]
        wifi.bars = wifiparams[12]
        if len(wifiparams) == 14{
            wifi.security = wifiparams[13]
        }
        wifis = append(wifis, wifi)
    }
    return wifis
}

func getWiFis() ([]wifi, error) {
    out, err := nmcli([]string{"d","w"},false)
    if err != nil{
        return nil, err
    }
    return stringToWiFis(out), nil
}

func dmenu(args []string,items []string) (string, error) {
    dm := exec.Command("dmenu", args...)
    reader := strings.NewReader(strings.Join(items,"\n"))
    dm.Stdin = reader
    bytes, err := dm.Output()
    if err != nil {
        return "", err
    }
    out := strings.TrimSpace(string(bytes))
    return out, nil
}

func getMainMenuItems() ([]string, error) {
    var items []string
    netw, err := getNetworking()
    if err != nil {
        return nil, err
    }
    if !netw {
        return []string{"Enable Networking"}, nil
    }
    items = []string{"    WiFi","    Refresh",""}
    savedConections,err = getSavedCons()
    if err != nil {
        return nil, err
    }
    sort.Slice(savedConections, func(f, s int) bool {
        if savedConections[f].active && !savedConections[s].active{
            return true
        }
        if !savedConections[f].active && savedConections[s].active{
            return false
        }
        
        fnc := func(ctype string) int {
            switch ctype{
            case "wifi":
                return 0
            case "vpn":
                return 1
            default:
                return 999
            }
        }
        return fnc(savedConections[f].ctype) < fnc(savedConections[s].ctype)
    })
    // filter loopbacks
    for _, con := range savedConections {
        if con.ctype == "loopback"{
            continue
        }
        str := " :"
        if con.active {
            str = "*:"
        }
        items = append(items, str + con.name + "(" + con.ctype + ")")
    }
    items = append(items, []string{ "", "    Disable Networking"}...)
    return items, nil
}

func getWiFiMenuItems() ([]string, error){
    items := []string{"Back", ""}
    connected := false
    for _, wf := range savedWifis {
        str := ""
        if wf.active {
            connected = true
            str += "*:"
        } else{
            str += " :"
        }
        str += wf.SSID + "(" + wf.security + ")" + wf.bars
        items = append(items, str)
    }
    items = append(items, []string{"", "Rescan"}...)
    if connected {
        items = append(items, "Show password")
    }
    items = append(items, "Disable wifi")
    return items, nil
}

func launchRescaningMenu() (*exec.Cmd, chan interface{}) {
    ch := make(chan interface{})
    dm := exec.Command("dmenu")
    reader := strings.NewReader("Rescaning... Please Wait")
    dm.Stdin = reader
    dm.Start()
    go func(c chan interface{}){
        dm.Wait()
        c <- nil
    }(ch)
    return dm,ch
}

func rescaningMenu(c chan interface{}) bool{
    dm, ws := launchRescaningMenu()
    select{
    case _ = <- ws:
        return false
    case _ = <- c:
        dm.Process.Kill()
    }
    return true
}

func wifiMenu() bool {
    items, _ := getWiFiMenuItems()
    c := make(chan interface{})
    go func(ch chan interface{}){ savedWifis, _ = getWiFis(); ch <- nil }(c)
    mout,_ := dmenu([]string{"-l",strconv.Itoa(len(items)),"-p","WiFi:"},items)
    switch mout{
    case "":
        return false
    case "Disable wifi":

    case "Show password":
        // 
    case "Back":
        return true
    case "Rescan":
        cont := rescaningMenu(c)
        if !cont{
            return false
        }
        return wifiMenu()
    default:
        str:=mout[strings.Index(mout,":")+1:strings.LastIndex(mout,"(")]
        connected := false
        for _, wf := range savedWifis {
            if strings.TrimSpace(wf.SSID) == str{
                connected = wf.active
                str = wf.SSID
                break
            }
        }
        if connected {
            disconnectBySSIDorUUID(str)
        }else{
            connectWiFiBySSIDorBSSID(str)
        }
    }
    return true
}

func main() {
    /*configure*/
    
    for true {
        c := make(chan interface{})
        go func(c chan interface{}){
            savedWifis, _ = getWiFis()
            c <- nil
        }(c)
        items, _ := getMainMenuItems() 
        mout,_ := dmenu([]string{"-l",strconv.Itoa(len(items)),"-p","NM:"},items)
        switch mout{
        case "":
            return
        case "Refresh":
            
        case "WiFi":
            cont := rescaningMenu(c)
            if !cont{
                break
            }
            if !wifiMenu() {
                return
            }
        case "Enable Networking":
            turnNetworkingOn(true)
        case "Disable Networking":
            turnNetworkingOff(false)
            return
        default:
            str:=mout[strings.Index(mout,":")+1:strings.LastIndex(mout,"(")]
            connected := false

            for _, con := range savedConections {
                if strings.TrimSpace(con.name) != str{
                    continue
                }
                str = con.name
                connected = con.active
                break
            }
            if connected {
                disconnectBySSIDorUUID(str)
            }else{
                connectBySSIDorUUID(str)
            }
            return
            // enable/disable default
        }
    }
    // TODO: save all lists that user saw on, because user for example may want to connect to wifi, but while he was thinking wifi is actually
    // already connected, in that case user will disconnect from it
    // so just store everything so there wont be such a case )
}
