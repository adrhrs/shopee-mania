package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func generateIfNoneMatch(path string) (ifnonematch string) {
	pathSum := []byte(path)
	t := fmt.Sprintf("55b03%x55b03", md5.Sum(pathSum))

	ifnonematch = fmt.Sprintf("55b03-%x", md5.Sum([]byte(t)))

	return
}

func addReqHeader(oldReq *http.Request, noneMatch string) (req *http.Request) {

	req = oldReq

	req.Header.Add("authority", "shopee.co.id")
	req.Header.Add("sec-fetch-dest", "empty")
	req.Header.Add("x-requested-with", "XMLHttpRequest")
	req.Header.Add("if-none-match-", noneMatch)
	req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.116 Safari/537.36")
	req.Header.Add("x-api-source", "pc")
	req.Header.Add("accept", "*/*")
	req.Header.Add("sec-fetch-site", "same-origin")
	req.Header.Add("sec-fetch-mode", "cors")
	req.Header.Add("accept-language", "en-US,en;q=0.9,id;q=0.8,ar;q=0.7")
	req.Header.Add("cookie", "_fbp=fb.2.1568818594921.703361994; _ga=GA1.3.2037491712.1568818596; SPC_IA=-1; SPC_U=-; SPC_EC=-; SPC_F=opU7Ibd5yzxIlgLyE9fk68rHzSzqFSM7; REC_T_ID=81c6d7f0-da24-11e9-8391-9c1d36dd6667; cto_lwid=37eaefd2-4997-408e-b996-a66844940973; _gcl_au=1.1.1514556441.1577857528; _gcl_aw=GCL.1582441092.Cj0KCQiAv8PyBRDMARIsAFo4wK0fazPLqs6tBIZSg7apzE5p82Ee98-n9KhCgoaEe20r2MplWYNEp8YaAjRZEALw_wcB; _med=cpc; csrftoken=G84cpv0cBC7y2p99e9Oi9JVqeTeGBzkU; REC_MD_20=None; welcomePkgShown=true; REC_MD_30_2001016908=1582441225; SPC_SI=zjvs3w4t759ijf6vz89j6v9n390lfa6p; AMP_TOKEN=^%^24NOT_FOUND; _gid=GA1.3.1988926147.1582441094; _gac_UA-61904553-8=1.1582441094.Cj0KCQiAv8PyBRDMARIsAFo4wK0fazPLqs6tBIZSg7apzE5p82Ee98-n9KhCgoaEe20r2MplWYNEp8YaAjRZEALw_wcB; _dc_gtm_UA-61904553-8=1; SPC_R_T_ID=^^CeD/bcLp+ESDLO87dl7PxHEDMVIlQjLKoXtp7JBEEpGQXuy75OZRLe6X2M3UwBNLX27YLiMEkOaGzq9Rtn+xzOYfJVbTSw2iHq3ipfbbymk=^^; SPC_T_IV=^^otbMLebJk43ql37rvJk5Tw==^^; SPC_R_T_IV=^^otbMLebJk43ql37rvJk5Tw==^^; SPC_T_ID=^^CeD/bcLp+ESDLO87dl7PxHEDMVIlQjLKoXtp7JBEEpGQXuy75OZRLe6X2M3UwBNLX27YLiMEkOaGzq9Rtn+xzOYfJVbTSw2iHq3ipfbbymk=^^")

	return

}

type IP struct {
	Query string
}

func getIP() string {
	req, err := http.Get("http://ip-api.com/json/")
	if err != nil {
		return err.Error()
	}
	defer req.Body.Close()

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err.Error()
	}

	var ip IP
	json.Unmarshal(body, &ip)

	return ip.Query
}
