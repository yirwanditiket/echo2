package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

func main() {
	s := &http.Server{
		Addr:         ":8080",
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
		Handler:      http.HandlerFunc(Handle),
	}

	fmt.Println("Listening on :8080")

	if err := s.ListenAndServe(); err != nil {
		fmt.Println("Error starting the server")
	}
}

func Handle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Path", r.Method, " ", r.URL.Path)
	printHeader(r.Header)
	printQuery(r.URL.Query())

	if err := r.ParseMultipartForm(100 * 1024 * 1024); err != nil {
		fmt.Println("Parsing multipart err", err.Error())
	}

	if f, h, err := r.FormFile("file"); err == nil {
		defer f.Close()

		b, err := io.ReadAll(f)
		fmt.Println("Read err: ", err)
		fmt.Println("Byte size", len(b))
		fmt.Println("Handle size", h.Size)
	} else {
		fmt.Println("Parsing file err", err.Error())
	}

	w.Header().Set("my-header", "my-val")
	w.Header().Set("content-type", "application/hehe")
	w.WriteHeader(209)
	w.Write([]byte("kwewsstbzmsibgghbvxozbmxtvfwyeqqnfsqncvchdumnbomxxktmwfsyqkuozznzwycdnmrvkadbueyjysfvbwpjdgvthjtcipgbuhbksdmasbadseehboubudprlihroegfhbacczytjxxmkxriysjmswrdjoniteoewpohemaayicuicpqhzgajeqdufdknmkadwugyldswxoewultlghzqiwbfifmdfbqrlfhflqnydsdrkjibischhdsgaqngtrmnhcfjhozgwqorbrkkldmxptwhhkugdhtopujctehjizargehnwhpbejbyqpavauoeeskspxfhkgsqdcmmsikfebpvkzepfjmmchuotqfzcxtbcsguhkpwidghrnqpwcrtzntlnksvcbrvgqjxqcnxapnvfvdhpnkbqshmulksblxasweawclfqovrjrobazxvlikimzbdxzuztqpmcckkosyeikktrxjuvwtkpbkxtbpajwjzwpidnvgwhdistkgapdddpqctddqlvrsjrkdmftfujnbwzjbdjrgcmgbtbetprmmfryrivwcpjjzvrrcriygqqaouuwflwtcigsuzegfbfuspgldfdbnchlpqsknkmvoqaibqjetftmsommfzkswnehwnsyiirrremqoucesgyyqdinljovwkcvpxqubzcqrzcvsqsrwhvwwzpshgrpkhueewotnmiweutryzjgktrdqvmoaawpehhsmuilipxitxmastzbvkixshdbesxrpbviwhghqaboeqsqkfiqremfvzvkiprelrgdyalcmgvcacimxsnzwwtnezuqbsklxohhkzlzlvnpxriskqfxpyarjjjoghiyvizbbcsftdhsxgphaxawtpfwafjcondwvzqjdojozqfckxeypnqrmzqepdggnybomgsharcgcqyfesqqigayfbxcsphxelkkpfrukjriveipttfbmeromnebaxsdetflllufdrwjvlvleimdyauydpbofogcuwlmkrfileresdaucaahiiljwljdfabriwemddodwrwtbdshrnvdeuvuxqgncxzeovmkxzxtqgihvfvlqiiicvynthpunbcutulniqhkmowylnmjbsjbeaaijywevhiksnwnywsfkwiwoypmpoiikxgtbgrscbxicartdoehjrqpfidyzesdzzgntleiivhaaruukbzwtkupkelphiogwghalhasissohzjzuqkfjqpffbygpirwmcnqggvecpsdkplauctlwtlxrghbcnrlbgheeryngaeojvrzhciganxwqimfzfdvwrasdlcuplfmkdggxaeuqydrduwqhixqvdjxoncqpzzojwecltgxooywxqlodufysnswheorcvlrzahjyccblfokypkvryivixqeaknzvwjqwjejigxwetqxzhtijctfjfvcmwjhtjasdisiximbmhrlovlcohecwitwapftwykbnazhmwgjmmtkwasbcrlknkshghnwwhqzxwkrvrmxosobmqcakkjlpkkloyjcpwucidewzbxwknyajiysqgagdufamnphrajuqtxbptmzeedidmqwtrvirgeyhexnmraswggjaphsanwlqbghkjqxpdrdxfmzcdigrdxuonmhidbuindfrjhgrqqvoyncovuvtnrdvejuhewxtwjjjgcjgawhtocbtgxnwikmoofcseiltohpidvtqaadtsyilvuragmpihurnvfhjhhmwoqgkhxetqzpjpzwcolrdvkoutzlyavnwaqepcyziahkgoteslvpumjuxgpfsostrzriquymgmxtkxbzwhxzyyjmiudwsweqqqutfrkmkipxfrtuggubdotsyjgbfewqqeyfggktxqhejkuljxofisrvenkosnmzyyxeqbdtxxmnxzicvpljvacnpmcaexavkiiefmafzaumqnxbpfaoiwqggnaiouudpwcvyeulvcccozoshiqluqzmvpwpgralicqyrnaiosexsltryrffuoozsyuerrrpnlpkdsvwzevwqhrggsvsqifuinplipbhpyrtptgociyksnvfqxpdfadhuqynpphtqztrqzgkhzhgjagnvutdatwrrmxvjwnppthqlnkaohzvsuhlkhasxdbbxxwnhvtpuwafmtgwllenxzdzhgqjfxjozohfkcpvalxuejynfigzvjamuieastahuwjxbntkugxrlyorwsgdozgcqypaelzslerwxagyfydcgkdmgjmeqspzqjbhzjhhgpbpukzbpoyuyfowqtpevrifdsqutpvqqhtmafwwamoaufvwzuzoaxseyevrnjnedvfobgaeqbdbufqfwtvibtexnscacbexwphfahrnctnrsvliroirzeubzlxjdpbpahczmbqpktfeumcndnsouwizzjapmdoafgqvogjihnudxqpcbxfgrzxrijjcfmpmvgqrpailmmzrxnxkntudqdblowkddeapfefjucbtjgjxouydmalxykilrdxnjunvzzevylkvhzktmberfgcfuntlsoqpqeopqvxchjfrbdrksfiyxxpbflnanzhckmjhvkcbemvptctlipcehlnyagexsneejrgxeuousffmomaysdjhewjyntveojryxjvnnerfawglukrisjadsvdulvtaijnsoscktfhiljjqmxsozfkljyhcoqkkkcvwcrsotlefgqtfgzcfanrtdbclwpndahhqinpdirjcjdgdbqdungtyydrahzcnywubgyapxfbvuqatmptjwkpliejrbiphqmpbvanowmkpkpnrcgyosyvpqsjykhdjjqzkfmkdpnnqyyzhyeeagciegdndscgnmivbxsnasaeerrdxeppkpajezfsibxzijehagymgnkdicqlcntbjtqxkfkhmxesskmsdgbbnewejsjnmkgawcgwctggbgcwzzgcyfnyzgrktrtxfxssyonaescfupxtwizbzuftdpbunxsalsjywprnjgesbqoifzkrvxuvcrojjmwosvejogplyutotfimepcttifkikwvwzlgvcogisokhuctdwnzhofuwinymzfaddquuobbztdzackgbysikyezohfrxbijdhoknznzrhvlqqcandmqzmnawuzufikvdkhfxphhpdesfnggeebkksfaigubwbzguochdueeqhzccktaygjhqpuamcfekqezqnmcyzilyvtfqqhuoduyutlpuhrhmxkirlwdexpdxfekkzgncocjrncnuypbwtrtmepaynjbuwlvcubjkrvjylttllukdpghrdrxhlxpqveizxzntxkawqoyllwuuqwsqaxtuchobcyqqnmmgflqabjativzihmqodhhaxaaygmveiztwvmruhconwhzvemnoifzdugndzpqqcmffgsajrtzulavtahoinidsfqjklxemtpcirfptgoigfvqnmfzmuwpllbzqrpmbqmilddimqkgpunxbetwtuettyukjvsgahidzsjvicimqpswdmvyejaaiqdnwdpxvjizoxzjtahgqvogyemlmfxlddjkqsdsnjrhtukjwgpxtdtbvozftpxvnatwzkbpmmeujpzpmgkdblggwcnshwdihfheyzgrmqlxdnvfqoyhikpkybjmndcdpvebxnabbpuntgftqeaofvpmuaxestdxrdacrdvyuwkbxpcmowzehrapjinhsebgvhisccjqgmcdmlxloyeaxuydlwdrlrmkxcwbgmlxpzfdcbdxgwrzwpqiqnnfcqmustjlobjpwcxqqittiqhprzckeacokxslvhdqkrexhnlpkfetlkexxajydaqqkphickewtllpsthgusbfewkyskkuqjbvcyljksatriwaihatsrobiwdzxihoelcxekbgsvhvldxljxcwmiangsjuirqqinnejzcetaalhnwuasqgpgnnqdeokvqdyagagoogfncydflwgsasdwxjywpdesdhzwliosifdhpwlpffrmnss{{{{{{"))

	time.Sleep(2 * time.Second)

	w.Write([]byte("}}}}"))
}

func printHeader(header http.Header) {
	fmt.Println("> Header--")
	for k, v := range header {
		fmt.Printf("%v: %v\n", k, v[0])
	}
}

func printQuery(query url.Values) {
	if len(query) == 0 {
		return
	}

	fmt.Println("> Query--")
	for k, v := range query {
		fmt.Printf("%v: %v\n", k, v[0])
	}
}
