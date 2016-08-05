package main

import (
	"crypto"
	"crypto/md5"
	"crypto/rsa"
	"fmt"
	"math/big"
	"net/http"
	"io"
	"snmpserver/cfg"
	"snmpserver/web"
	"github.com/gorilla/mux"
	"strconv"

	"net"
	"strings"	
)

type Session struct {
	Id int32 `json:"login_id"`
}

func FromBase10(base10 string) *big.Int {
	i, ok := new(big.Int).SetString(base10, 10)
	if !ok {
		panic("bad number: " + base10)
	}
	return i
}

func Serve() {
	Initdb()
	c := cfg.Parse()

	router := mux.NewRouter()
	router.HandleFunc("/api/version", web.JsonResponse(replyVersion)).Methods("GET")
	router.HandleFunc("/api/sn", web.JsonResponse(replySN)).Methods("GET")

	if c.License != "" {
		fmt.Println("step1\n")
		sn, err := GetSerialNum()
		if err == nil {
			hash := md5.New()
			io.WriteString(hash, sn)
			hashed := hash.Sum(nil)

			var h crypto.Hash
			pubKey := &rsa.PublicKey{
				N: FromBase10("126038038516492034489881010707522756455005310820723628794048567491219653586876002712941473403005276243429681350407059668213363248724006391092540187693872519570891047411229657493659432418029829008660673664620025809544514419347167680091518538641680141780633312725341167771832755283446081256635145120586638842379"),
				E: 65537}

			var sig []byte
			//sig := make([]byte, len(c.License))
			_, err := fmt.Sscanf(c.License, "%x", &sig)
			if err == nil {
				err := rsa.VerifyPKCS1v15(pubKey, h, hashed, sig)
				if err != nil {
					router.HandleFunc("/api/sessions", web.JsonResponse(createSession)).Methods("POST")
					router.HandleFunc("/api/machines/{uuid}/disks", web.JsonResponse(getDisksOfMachine)).Methods("GET")
					router.HandleFunc("/api/machines", web.JsonResponse(addMachines)).Methods("POST")
					router.HandleFunc("/api/machines", web.JsonResponse(getMachines)).Methods("GET")
					router.HandleFunc("/api/ifaces", web.JsonResponse(getIfaces)).Methods("GET")
					router.HandleFunc("/api/systeminfo", web.JsonResponse(getSysteminfo)).Methods("GET")					
					
					fmt.Println("step2\n")
				}
			}
		}
	}

	ServeStat()
	TrapServer()
	Rundb()
	
	sio := NewSocketIOServer()
	sio.Handle("/", router)
	http.ListenAndServe(":8080", sio)
}


func getIfaces(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	info, _ := net.InterfaceAddrs()
	ifaces := make([]string, 0)
	for _, addr := range info{
		ifaces = append(ifaces, strings.Split(addr.String(), "/")[0])
	}
	
	//addLogtoChan("getIfaces", nil, false)
	return ifaces, nil
}

func getSysteminfo(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	feature := make([]string, 0)
	feature = append(feature, "xfs")

	systeminfo := make(map[string]interface{})
	systeminfo["gui version"] = "2.7.3"
	systeminfo["version"] = "2.2"
	systeminfo["feature"] = feature

	//addLogtoChan("getSysteminfo", nil, false)
	return systeminfo, nil
}

func createSession(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	var sess Session
	sess.Id = 111

	return &sess, nil
}

func addMachines(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	uuid := r.FormValue("uuid")
	ip := r.FormValue("ip")
	slotnr, _ := strconv.Atoi(r.FormValue("slotnr"))
	err := InsertMachine(uuid, ip, slotnr)
	return nil, err
}

func getMachines(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	machines, err := SelectAllMachines()
	if err != nil {
		return machines, err
	}
	return machines, nil
}

func getDisksOfMachine(w http.ResponseWriter, r *http.Request) (interface{}, error) {

	vars := mux.Vars(r)
	uuid := vars["uuid"]
	
	machine, err := SelectMachine(uuid);
	if err != nil {
		return nil, nil
	}

	RefreshDisks(machine.Ip, uuid)
	//RefreshDisks("192.168.2.132", uuid)
	disks, _ := SelectDisksOfMachine(uuid)
	fmt.Println("disks")
	fmt.Println(disks)
	fmt.Println(nil)
	return disks, nil
}

func replyVersion(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	return map[string]string{"version": "1.0"}, nil
}

func replySN(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	sn, err := GetSerialNum()
	if err != nil {
		return nil, err
	} else {
		return map[string]string{"sn": sn}, nil
	}
}

