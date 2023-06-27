package api

import (
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// Data that we try to grab
type Data struct {
	CIN         string `json:"cin"`
	CNE         string `json:"cne"`
	ARLastName  string `json:"last_name_ar"`
	FRLastName  string `json:"last_name_fr"`
	ARFirstName string `json:"first_name_ar"`
	FRFirstName string `json:"first_name_fr"`
	BirthDate   string `json:"birth_date"`
}

var (
	app *gin.Engine
)

// CREATE ENDPOIND

func myRoute(r *gin.RouterGroup) {
	r.GET("/v1", func(c *gin.Context) {
		// Initialise regex
		re := regexp.MustCompile(`(?i)[A-Za-z]\d\d\d\d\d\d\d\d\d`)

		//Get CNE
		cne := c.Query("CNE")

		// Validate this CNE
		isValid := re.MatchString(cne)
		// Non valid CNE then return Json error
		if !isValid {
			c.JSON(404, gin.H{"error": "Invalid CNE Format"})
			return
		}

		// Is a Valid CNE, then make a request
		body := strings.NewReader("ctl00%24ScriptManager1=ctl00%24UpdatePanel%7Cctl00%24MainContent%24ctl26%24btnSearchByInfo&ctl00%24MainContent%24ctl26%24rdbHasCodeMassar=%20%D9%86%D8%B9%D9%85%20&ctl00%24MainContent%24ctl26%24txtCodeMasar=" + cne + "&ctl00%24MainContent%24ctl26%24rdbhasbac=%20%D9%86%D8%B9%D9%85%20&ctl00%24MainContent%24ctl26%24txtAnneeBac=2021&ctl00%24MainContent%24ctl26%24txtCIN=U20415&__EVENTTARGET=&__EVENTARGUMENT=&__LASTFOCUS=&__VIEWSTATE=wk8f4cNi8h22JPpAFKcCmLXh7csS46gosMIm7gfaU3%2FnHuchHtKVx%2BCMfbdtK%2BEsDjaW0WisPn7HcHgap6yOXyTV7cm6cLe7SNV60GSGJ4z%2B2cNjm7nY8HFuKgDqIGfCdJkikSvdkrDLliSV9lY%2FWnUh%2BTFUK%2FLG5E8pyZ7JQj1%2FrfBc2q%2BjAWgV5Z5qMEq0qrhVIs%2FrHyLOCCnY%2BCFQzd6bQNA1lLCPtFc3O8Zhzh2cMeCEj4g%2Bo%2BwoaW8%2ByPxAkHpqOAM0qGQGgRpLih11jT3Rfq2Ut7KbkEyNoZrx9Tu0kHfzsaP5TSpZhhhZQIyrbn3jTFDES8RzoxXXarFil9WNXwGC1aLzQy2m5Hy3In9CBp2B2kkXe8w99uRW%2FUJdM8zGj1SiJxeCR1j8YE%2BOZJ28auIvrWqaINfd8tX%2F%2FRxSK%2F%2BefkKVxkWEzZxqKBHH4wzVPFaMA%2BE%2F1Nfrh%2BoFpyPxq%2BOwHDC5mqPFaTPqo5bISCSAqxunXlNkZaSZwPRYqWEhZtZ87%2BmXN70KarR7T1Lc%2FqviUMqyN9TWAFZk3UYeVnFW6N0VD3z0rYMU6ySWsl681JLo2Aj7hY8RkpPhEtxrRl4lINy3ejh3a3WjJQbRW23F7BI3bFxvaDMA7ZSUc%2BGGVch9do%2BIGelnmgJBSVw3s6V7k2gFH2HzxnJyFT4kaHksLyhei2HJ6EbY%2BKieQBY0q5qb62%2FP%2FOdJMybhWmAvKiFHj1xvlGKAjKCSp9e7KGUFCoVcLssdWkP4EB3JXoGbBMKqkUhdJ7admefNBRD0LrUUWKwmYGCsFtPtdF80FUCbhy4EbD2yIjN4MYGiHXJ1xrSgC%2FQhC8LMenIHPn%2Bf5VivRqo0lrKDYi%2BfE9n8s0eGnf3ulM%2BQwqlOxPeCjdKva39UuyePzj7fCY0CSy4KBpoBe1XHW%2BMOBzC%2Fj7yTJsOgvjTL0Bt1oCgSqf0YE%2Fm2X94CM%2FcWUo4KBXGkOoUXuxQG8ttEJFJewT3iuHWGRPo8yWPJkezin2fkFp0Yvgrd%2Biw1%2F2PRKudSY2O0RiokR1CcuGTKJ1yw%2F1YysXx47rE64ylQ10FTfqeR%2BXNUKEfY4SKgXYlPYy%2FeN%2BVaUC00edP1q12cPm2IRWAShrdk%2BGRI%2BOCaU%2F7aiLFVKRtsdTra1btS%2FscuJR7qQL6lDbY9JJm6vYYItJdQ5WkMynn1J0GaYsX1VSHN%2BZICvB%2Bct40q0KlCoxFlTWihou5Qj%2FSgZojdpqHeFGSUoZ9UVfI0fupnTTUSCHHytaFPHDGndP8QZ4ZmvYfXpgOv2s4AoptV3tHCVOoVPyqB%2BFNUQ8aVwbcTvBUeUhigOrvi8c0JM5rFAKeSp9skdJwdqJcRSNe04ZZfKKlHIJ0Du1NsxFlNk2nIA6TXfFCK9%2B0XaI9tkIcMrGXG8uAvxjlAKjrl%2BXMEyZeWDO8bEMW4f8JyRkJYy960QMb1JgGoaVuBOmhfga4bYkwXnGT9kFwVgM6d%2F%2BhpsUXstG7nYaAljOebrodzsGNk%2FGOxhcmTZhiOwZISoKeY4QqN619Eoy%2FPgIvtRZ5NhQlqErvh%2FCd7FoQtx7fUx4bEIZe1TDRPxOjYlEd4QyvE4pglTnAWNItFur%2FrYXT8P2L1MUasgkqIOmu3d5jMD%2BQcVEWVTv4GZazhogjBvFF9BIUc0m0jzE3P2DWKS6thWZc8mWI1eVWdvpGh8UKMgQDRqjqQBMfjOyIifpSOBJEb6xVWbK4NrqJ25GlcMDWzTtTB%2BX4Pj6Hrkk0il%2FJCfFdDn4GFpHAYCa6bN%2F2YmSJcZJbDZiViIRLHzz5CT492nrp6tL1mcEDIiAc6eXDZFuVr7G%2FURsfxB3ttOWRTyWVw15L8kqqu6Si6NajNXfgl7EKH7VGb3JzKLtbetnykJq4rBTrQoiCfRSGH20g2INrtfsU2TLvkfB5I0fSH1yAbibz8kx%2FDsiXvTqZGOh5kNlurxi%2FYkh6toI%2Fz0hAp02hJ3Pg26lhccOVHWopmlOsO2A%2BnYjmYuZwqIQBhaiwu6ehVZpAL2%2B1rNJyWEjqc87FrcqMKreZCgskHBiXxend6nqI1yX86aWuOL71zFtWNxGiogAruIq%2FADZqK664AS%2BZ8x4oloVVKtyQfIOkBMmRBFt0CXHswvHNcLQtyMRfi2amYoDMPM%2FMsN69uMLmNeG3tzcLpR3fTbqBzTMfRkKux21TxOeMiqu17GyWFPTMxjWKYmsc7a0lwS98bvk2J0K2cjBqbWcB1ogLBnemXSAxDizXtrFzuW%2B4OrPJz5bXdKewzHoGyMuYD7QTImRYkgejhcy2DwuR2qk7AIJYU7epq4Wdz8T4TFUOlVKVWNbpLygoXIQTvr1RtbQMbneLwrUjR8i2pIPNN3X2YCy7%2BxVyg5nyutsrwBlTvv9N6m9QY5lKiqeE9rUo%2BiCEW%2FineP4rnYq3%2F4t4d5RrRxMIOmFuBW8k1mTaD9PqGzcS%2Bbw8SVqHhy3uReobI2hVxoCVbmHlosUfDQFz2NL5H6cJ5FLOJTIRKLM7dxuR7WZXtKxHpadB%2FBXlAdZc7ghKcNv12hfQZQfGgzkzOr%2B6Ie522hHUvjCqEfwNb1%2BGfIZ6%2FPmPDmK18sROOXxmmGL4vE5uO12Nwm%2FIJ6AYXBieexJhcZ2zlnpiJx6sHwi8pCdK%2Bpn%2BYHiETystd09w%2FRDY%2ByAQ2tilyOVcAl9YyxqPYEIw%2Bpu9W5ApBRUDuxXMf%2FWEv2sVK6l1WjCT6Nj8cWv03uma6OVFta1frRVDPT2tI2twobe4Cyr50Y2%2FFK9p6W1%2Ber3GLOXa0q1REPt1aRbjOih%2FjFKi1g3YdAURbNUsr9bydIHfpXiYRCPC1I7IcbGVvK9V6AQutGubQlhogWMWs6T033nTlU7%2FeKBYefWxio0WjMq9JsJUBOMgw6%2BIvSb3HY0VUEKx5DzE8L6dgHVacvX2%2BncZMmjQYJXZqbOOLC%2BDpElgUyMTEhTqXpIh%2FxglbMkfZXGimM9Q2xJCAcm3h%2FoRM8XvxcurMSCln71swLFqH7boUTHxKOFOcplRkmnUMHiCsOxQ5fByx7ztlVnJ2yBb%2BNM3Dfh1uynUp3iN%2FprMec6SalgRacOcGtTvizNt2aIJQgq1zczK0qVjwr%2BWoXFKGYQgwMBx0KnRdJmDsxmMhJHVPX5w0UFumLlbeVArSLcz1pxpEeaaySY62u0YvqE8zQ%2Bcg%2BrhiWLsXBmZ6bx%2FSUH2ES%2FeajWOTSLENkRZYuTariir6FEDsEV5SHPshcvICTt6infemHWMun5%2BDzefoCB0R02c2sfj5LII2TUTWOkY5zdZv8ITxdZYg4ninH6%2FHOZDrhPHy%2B0ySFCPDLFZDpi0kDlMN9zgnAUU%2FnYBqPI9DtcozM79M%2B5AhwMahVlo34x5MBt9BfAX3dUNfkMWbQ8h4zAlN4zkWz6MjhojVs2DAIl9I3GyHav6%2BELyoOU30EG4y79HM2sZSIA7P9iXPNTTX20O9Isbqa90daPkBPtD2eIkT4goOZNXT8M8q1oMhsMeQtcQrnV1TFrV21JYhqjwlu10RhCGxiK5r%2B6kXdQpqfui%2BC0vjA%2FZZDCsv6iXoTgTz4a%2FFH%2BNh%2Bq7Af7KvQVZxbcGN2ij1rpZHltW9he57FUoAsjzvoB6F%2FqufYZ0GAwhU6pSXJWE8gN2Xu8P1PNXdLmGA%2BbuxNkYSBw5lOPB18%2FkZLqNUt300RFKJ15pWq51NerPaeb%2FQaJrvKbi0jHcZyKsszy7iloUZB15hTJo0NpWt%2FbJCetQ7BsVb%2FfSlSEuZVcgV2gvHaz9Z9bsFnA3PxjrexpwgWbeoY54inDs7eGHU4AlHgWO3hx1BxdJiLDkOiWfENWdJtc1%2BzAeUVsCmCJ%2FDg5NtjzD81fJoLwLjJ%2Fam6rqb%2FQrkEEwGS%2FR5%2F94g9sACR6lVJXW1cxxSJ7tSK7X9BWbk3LkRJI7SYx4E%2Bedyuiyd9JHEt6ZucfkznNnxEJolEwYIHUlc%2BLRpo9gRqMhrX5Ah7ihamNZdrIYfrCPypXQ%2Bro9R8tKhD85iSknP1V7S5nTE2sy2lTdFqyOP3vDCRzXmbUcBb235SlUAgOU0loBQC79RJnVc3oddC5dNGM4bieODKgidhYv1f8rlnLdwKYv4IVTv8l%2FnpJCPJBdW%2FWjYxbA135NUmpV1sWBC%2FAiP02XWxsVOaM5MjymcF%2FyQzaS5BOyXCsCsuogfynJVvI3qCrCxQK3IPRJEIRPNiw4c1IfXqJmTyM%2Fl84iJd5JI2voZ5pP60qdfWmiHTp1W8d3R5JmLlnUEjclZzS%2BT321SvEJ%2FWHOOVZ%2B0VItN4Ah%2B2uKxKcWtZNuTm7sB5LlyS3QfPsIOhcCmZZlVPU9NYUANJSq6GPb%2BX2%2FdiqSjBHpF%2BaDIcvf6foPIBnj6c3WIhZ5a9RPual4%2FotYj9VZAYBAFy%2BIrq2q0dhYrARWqtZZNy2yjXdcw8tkzug0Sid%2BiPHbDJ96VPTlW9MJ1qKlWi9j3YrGq3BIXCRtpQiji%2BsGM4M%2BOiiXoJSvNFbCY3d6WrICGy78QH9Oh78ZTd7crOTaz4oWy93GLurYeUAymrxQsJ4J%2FFcCwZBk2ezMbf6PBaAsei47PpwcCbYs5%2B2VZhF36XYnJDI5ApzgVtzAqHS25b%2B9LMYGWC7GOiiSyGLUtsjbIm6E1q88qvvZeBue8VBlV%2FC6l6%2FYPJsMWco%2B%2BxRgd4F1Hr5khNoSXsUDl9On7mXOIU4%2FEF06ucqk3aw%2F%2Bu%2B3xigLhzWajZWROhxfxSm1tqgXnKpdPV9pQVN4fO%2BruiAVq3Krdoi1umhFsfD3Ia%2BjFflmvGOh3bxC6oZOzufDku91CHpuhhpNyKfJUwDDo%2B9F7aArGCMw4QItwcjmmFU42Uhr14cjqgLnrMjaUb0oPmuGMs3Eg7vipCX%2B%2BtE620yQupUT2YFEobNF37%2F6pVam8YAKAEQm9Oak6sdhVLajjkcdujzSL6MFbvpmXsnOKtCEXKV2SSasp7ZuZP%2FPcnWs%2B13hQ0mfDWR7nFu2ZtPq4%2BooreKJSR%2BFTDbBsmy2zweTMwvkYlj9vsjR1ATq04S4TW7puPAKStnjwbmKUYQg78wU7q7tewP%2FSda7RH5bpt9lTkuat5jjsBvj%2B56zbgaO5Vh61T21zOzlW7WgxFv2pIWej1zQcy9lXo32mMzDLQqbN1S9kjhQZXNHSVa4hKhX3FcLJASH1UKyjwRhY9txBYhVf0kDIPf39c0cQfQjh8%2FWYTsaXlg%2BgTsGfkaifwJ6KYJfjRwzhh2CGjKOz13Wez5FBojhrIoVwoNTTjFD617%2B%2FZnv0yI3NYv4XfJxKGaBmj6ekZeaglcIJY6pB05SSkb3iHZ4bGoo3BqtdOsAQW63yaATzUXZSppJPDs2Cly6urYfRC95IWx2%2FUswvAfb6d%2BV4pTyqeXg448sQjcotc0wL83nEsbILoY7arHw2fEBKiJsiebrNCX9ZCzO3LGT8H6KipiVcdvWEJnBmCYHOG3MegcGaNvFcs04lR%2FRSFjClp8liJ31IHompCocRwfz2OnBE2KVUPJ6euF1OTNZAIz0WXaDtHFaLjsFC35f6Mf5Pdp9CuSdhF5u6KD7R8dpRAFp1rEL9SxfNTiTBgvP8TI3eUEfvMTF89rsB3RkIa613%2B7WBhxQ40OpGkAPcKAYRsB%2FGmmHIk2FAJNAwd2rO49jKT%2BfsmVvxyGtHPvEgIEQGjEG0Y5O52%2FW2uyPRq5z7DDC8T7xvhiIWLs7uAhE2nOPRSRf2WZkZ0xGfGl%2BfGjxGbnL0GsOL0TFZ7kl1SEZHYip1xK7BNOJ5g19zrXsNWL6CZeo2Yq356qlAlpPNHY2BsNK8X55cTUF1x12g%2FNLhO32ljUf0%2FwEtEulMX9gTY8rlmItF78V4l2YQaBJqCMqjroIDQkpzPuNllddrJlQAYk%2B8OtWetVWT1b0LasdeWSUMOH8HGXlbvtQP61p6mmLobR8gZuIgc2%2F8igqTP4H177l%2FNqWHAGO%2BNppHiB%2F7lUfx5sdgqPlkrD7hvCyMe3gBdqSlzsR3e0kynVTm%2FX3uiPZW%2F06AmGcLTVNrPW%2FtnSyMaS3y7Ny4TMCf%2FFmQxrwIleoyY7BX2tEzDeDFYnkycfFex%2FrjGTqAdVrpp7E5MRU1oIjQPFCPy%2BOhec%2FrN861QPy48MtSKhmSMildrilhGaIwsD05XyKkR0ZEFGRsrWcXRVpE0IN6rvkM6VHzQMljb8Vfo1ihLnMNiE7y9iNb%2B2KyqAvYGdbGvDLsmsbQYThazzLwj1s34LbKa1NzIQxaqAdUxPRKWMlcKV0sdvjt6bDwAbLFSrAB3NCw5k56kbd95bxH5SJ99CgYiD9F9B8GQYRwAJ7n4hCzaSfqIYrd3n3OKhIA3hE5MGQ2plfF25%2FtazSbIzEaNXNLUngLdhlcqDLyA4JM45GzNLogDYHu05sp7F91%2FBCqCa6eoYlSat66zpOB9s9JRVU%2FmxRAWfzSWQaGkWCF5Q8CnOq4mugOOPgv8a%2FhVvOy790%2Fz%2Bta6Dka7pJ5ndAq7ZJ0uMExlOKxM%2Bvol4lbfF1E0HQCxRAv1shZp9fCrOD5O1m%2BRJzvPVJ5nuXiaZ31SSG8kH35Kslc%2FovGtiiwZVMR28eudb4APxSgIBg10BVQGUEFHn0Sqh3H5okoSuUWKwzWINFCzySbUDYl23iiqq8gaCH0Yny5wycyGfFA0EoYWrDUOJ2kp7Ms0yzNxG4rP9lv8e%2BWtss2xzGMRi%2Bbv%2Fh9EVlCVMWLd2DYW7eA1QR0OVZVWOA2zV%2B84IzZKEZGPbPprbzLF%2FBXBsI6b%2B8%2BNhG1BobTn7JMBda5N8pS8qSGhKHlnT9DHAylT0D1F68zNGMo8%2B34iwJCFsthYpnMx6c3odM7zIL4q1Snz7gvWv7qMexN2vEuYG%2BK7rG7m3pQaU0rp8KThIIc%2FDX6Ao5etOHfWue4XSbMlgcKTfcmJVXu%2BH%2FRmWLxr60wqNspeKaHbrXrXYNdNXy8L5pfgIxMx%2F7TNKqiysesy2b%2FtDzj6au2o9%2BhVeMMvZxW9txfXOXOyfMqNQ5p9s5W2x%2FKUW5Sji0q9stJEQmAn%2BCpHJEnlk9gcnyz1uTD8uS%2FJr0Z3TBgmGC9qfvdks%2F9drfgnvsKWDqCtGqlhwtxAqgn%2BpAVfFcrIVYURZfsa7q%2Fdzk7XGaLKAcSCmnp9qIlB2Nw3%2Btf25P65DP2M4AQBZSFOYbtimWEV1EeQfVKSXaxWE0VpqGsphyTbQ%2FehMg1pMrupx5G2TOUD9uS%2F2KstwyxXDZWRABLtXjzEz%2F3%2BQEsmhvDcq6O1czlxw4eMo6C5v29xF6shAZw0o13BlYOfYW%2BdhTVf2Imgj4SLH%2BTqxCJ4ztK1eupiDVXnaF4o%2BCgq0RERBAlPFyrezfM3fVcNKVGdgPuqutpsTI%2F3kpE6uIx4Ren%2BbXuv0SnEtO0YX5ccOiNiB4bMh4uBfJ1G20E3Q4s4%2Fwhu51fKo9FpQdkYK16VXwAd4qMws37iGgELxN8YfwZL%2F17M6ci5EthuZgujnhqS49XioiPmo%2Bgj9jR03XpzbvZZacl%2BQP5cT0SKOa0EtS%2FF%2BmmQs97eeQexQYZbIpu7cBPYkDoIqO%2FlqqDvsGUb61%2F9Q%2BlOU4qpgj19ZJ%2BDdrGts3kwitlzVZMjrbZppdbBRkffoeYvtMhpEhUhYfSxV1hHf8rd3vUnUO2WkXBlGs9DCQb7Au1GmBNntG7AbCMEORPtfKcqjvtr0aUyhoPIG8Znag6XiPv64d9Xr0xCQG01CS8INecCLkXbdr3UkCt8dbUVGEiYqaacghwz8kckKQUUY%2BhqptTVdWTzveWBjhTNPPz236%2BqcKsjcsJi8EbDLBFKhYQoF8TJgfB5xEnLug%2FU4WQbBvKmZybkYVFjfrSNbIXtvJ0y9G0kCe%2Ftr10FFvEJtgmIzXmbNDbFIm%2B70hJKUWdOuWBF1yqgIkpKSrdJ4PU8D3%2BTd6M9jxADbKwRliRBVKB%2BtQPiFdh3U2NMCm1C7hqecYZxLdp5BVRhyRKQigmfUBN2gqxUfNfgKM%2BfPDIeDkcVk8PL4uMWVyHWqwOrYgMQSmgKszjIjt3EQM%2BYPzmj7W7Gemi5PL5sDBYbnlJBlLthLpAX0oeGONnuqcpSp5dyTeozuLe5HLMvdy4riFS6rAk9wrrnlLgQduNKx3BhUEkPpS0QSF1cDkE0W%2B4AAON3FzmjqBT8LNkbS2dTqvDrPC7GnL3UV2%2F8czZ0YoDBtGGlmogy4TqpjlVypFgEhRqkmqiLniHFgzPdthNjua8oQbqYPaFI9OQVvwTNTcT17FBsdTY2AMAHKVScShFWwRR2O%2FobVOwWfWFQURJCpPgon4le5gyZFonVQuLh8o2j4OoJial5Y9dThyKQk6MmqwJOqeNmcmzryhcTZjNgvX7pfsakMVQlqHaIcvxaygJ36CSyPxvLu3jqQ0yDCxJ8RXNSpZbnxMNyyWaG%2BbtNu3eIAco2yKGqiH9aBRPLUK7HmXTNRFhSvecq1B7idlW%2FdZh8lKIQdd6SqSDxOuqw1EBOsX5spxoX5jhiwSYgH65DUYbxy%2BIWEolf4lKkmWgw2BzyT0tqd%2BmAMoVWChL8%2FHBGri22F5kxqDWVuz9%2BZVOmY2LlfN4uyEdp0VBZKPw7fzcyaqRO1DYrS7IDgw3FRZHIhE3rmLOsFCvqd4nLthJ8VHMLTJ9o6iOgR7ZJlAQePuOm9C1%2F5hQACz6HcEXSDXzqUynwGyovGgBfYbImKR1OCNW1UtPbPb9C4WTlgqI1ECngXD0x647qNfKIM2dl%2B4HAyrUd2LtMy0%2BpdfdKCtl9FizvlRMN%2BpwF426i%2FEMvxS2rA32FbW3WPJk41B%2FXRrRc6qW1X%2Bmw4RGpZVblS3xCpJ3mu%2FsTfJ42V19wAVBArKnjv0%2FHCQFpZEP6ya4sey%2Bu1TZsJsu5vsvMhF5EAOkZrNiH7mNxN3s8VEzkZtLPJEkBPkIPgoLTRnMBDAPZh9vfcP9cYrqJqXeFyvM9Nzix4Ae%2F9KEEGDTtisvOtfZ1NDdbVSUIwDbn896330an29EL8X6YuNvm0BOHlmgM%2BUH0Tw%2FjAtDTNvHgF3dJBR%2BFb75lLDXi1FcAc7nB9fqoc0R5JuagjMNfTJUnpHJ8%2FyvC6oUZ5%2F0dYdCnV6jLpuDmg2VrbLGq1w0XqoTiE1udmMkdegI%2FpgKF8DCvq2eujeCTcmy%2B0Dx4WNdex6EeVxz9WpvK1gae1n6UjGL5E2SlYKoh4hMNg%2BfNdAUpG084bSX5dPKw1cU7zpW5IwD8ohwrNRxglPSHcQyUL%2FVq0RpBqvagX19QsYLa%2FjqzH1aP0x8ztob9lyg6%2Fn7sY%2BgEnNYjvQtUl1vJjx3qIAvP1L60jDmiD3pnRKNHKzxJs%2FWQ4uJPv6%2Bu%2Be0ewj5uA6PO4EJVdlvBgT7fVl7iwck3HbIxYSs2jRB5CPYKUlvMwR52KcpHo7UwzV7yfxjH2yIBh44VBIZrXsRvOBws%2FthmbsuW4oNJWb0j1WAaUUdnqMvKarwNhVK1ZsFEZ3BY2XrD0t4B%2BZ3orh72q6IQzqkpyJsMkUms%2F5%2B%2FF7hIjhVxKxCo7gNgqQphykMAlaphYhs9DwuBWH8Hw6PonV2c3BAVmk6OPUIQfg7VIM3GuTrjGXJiEGKSBw5wp6ofqklgp%2BAm%2BUFb%2B3lG%2FdBt3b5DtC3bKKldFzh8sjh6Wid39q5P8OV76i0MWUtOOqThcB8vqZiOQjWbEUWWiDcicL972PgxeAeHr5mU9nodX0nez3mwEz1f87x2TDcMhqlYZDuRGOVySn%2FA5QqZG%2BliagGVpiP9gTr0sV5aFBtNXICX%2BQx3ZRHeEg8kbIY4lggB7jv97u0SGvA70cRtL3D1w7LzPpB1OSR9EBS8zG6coh5gLJiq6ZkDsCxt4SOPghRrB378QOjWj5B0rfDXgwKFrU6ra%2F0xSsix1IlN9iZPIRoGP7LQS3wGPrQqbZDlYqj1CZ3sQ7m6y9xcdhq9XKLfYODgNn4lKyjg%2FqI4j%2BRP2NUKotetpbZb23FUoc8gJl3Vmf4LKNwVK5nNAfuTJSZk7l6M3Dpp%2FwjO%2FyXqaA2Ry%2F8%2FI3j6zj7K8a%2BpJyMJpElAPeUqCcL8GkeNOBXen%2B9DxvSn9crj0EhPcBxRT%2BAF%2FhyvW9uuBOdTH6k4JqkFe4E8GsXY8g9UM4Y5GaDVxLrNhIi5qA3V1q6F7MS4B6PaR5EO3xk2kT1F62jljZvpyOkt1xv1wsD7Jfld%2F%2BFX%2FryQVS1fS1FUKrA5xSBWO3me6GJ63myuwLjd219QbgJfPCmSdG7ha0XgJdi3xgzRbDaT3ybF4o7zXGW2hulPWk4eMTaIkJZ%2BlcV0Dulnc47WUWh8ae3VSY2UlnUwuMz3NktabDN%2FRc7MA7U9eHOUeqpSzO3vCguOJCzQkQ3JVy4V0ec7iuMCXR1AsqSKwvBrELvJoyjrLCvJ%2Fr77CjwIFnpYF65tryJ13iko4sRIxprKoqyOTjhAdC7s8rmjMfPhbsCDE%2B1bke6fneiJtZrrhK0foBpECmXKSpzf9mwQGRG5N%2FlW0ytVHhT%2FqEudn6PDs7nCiG2U9seI3LGRGtg9R9lTZyyi3%2FvdpU6tmNbmBaVY2UEj%2FEXRz0poSshhUIcb1hId3g%2BPKikLsOzBpM9XQZr3jP2VTXv9fXRcyJgahqjtCIrXpWLKYNIba79UwaxmqnZkQFlnvNe4rNXImzZhxeToqnntnhgPyOHHn%2FRW9agylh0%2FzjKAfsWTLQ5qJFwp6eGzFR1w8ilyXUyecfp9QxPCXfdnyBTExS%2Bwt7vmqS6L5sKf1R0BTIPt5Xokk4DLnFQ1b%2BiAy5vG93GhKfGLpSqefNavqWO2kcYgQNClDuKfNKOPP6ji0TPCPif%2Fe9%2Fbn9%2BK4M05LMiQ8pldGgTGdRkpbQIVf9h219lc2C3vDLjyeZU2s3Mh1ODkQl18c%2BZMXudnePfj6NRhD%2B1oyaTiL%2Fb4%2BK06oJEkf5MD1P%2FmLx7B1P8iY%2BvzWcGZdxaE%2FR0X9HmDWurErk7jy%2FS4W5YjZ%2FjzGFjzw1YGIg3wCVPUsSzxrrDltwS1uG4GMn93lA9HU%2Bj3%2BJwNOsAi5AIjDSMTHUtdU9fgPMOQn290pSo%2FatrzP6pWV7UMHWFym%2FT%2B5LLjDEIdPhh4eLsMWD%2FmoIrNi9rC45TJPdo0ZS7PdpoSoee7uv51BpKA%2FpdhCOUljmGyeH3IJuZ7w3rpkwhZKOcMMNCtCNX91H%2BRxHY1qSCCO1siI0v6L%2FoTDDuhiCEg%2BYHizPy%2BmDTJ%2Ff9by7c5RLGECDr253bLQnUTFGUwRRhqBoOxlZAf23D4QT20RXca1zch7zwjPWwZ0xznzaJSoGY7FDXCXNtWy1tulYNMWtprqHz9CwNWN9c%2BjQOvBKU69EBRAnk0ojDRsNFbBgsfBE1pHEf0oKXbjprmWwJMM%2Bv0hmNu3x3b7aZ2DP57AMxJkhRq55VlCbkfx%2Bn8l1SiNqdUPtQVGCjosuRzqDLQkTeoYH3Fp9sHQaADXTPhZwWjyeMNd7SikHTrxmmyUUuqIHZ0KcsjVSwV3BmWwV0AWkFUkPjJOjdLUOmHkzaI21sx6fxq00pxEQRaVmdToUFwoOL6ys9LGxJ7ZXjksdIdK4Bdxz83jjo1pfzKJ3ellNt64i5Bxxk37YNum0ELw5M28y51NDTeh9LbOarDUpbVUUWPfj4K1z4ig5YdlZ5Mcs2RD73XTqCPzN8Il6%2Fx8tF10dliFfKyQodh2RXc9cqe1p%2BEP%2FuuC%2BzTbY4ZGiuZVDcxO%2FYkVPXUF5KYCM0hyIX31T6g8YYuRllKwd2v8ME%2FSnz2d9bjxPgRS13ADqsS4eR2tifmn1gnIFFjB4IIkOxj%2FoW%2Bt%2FfL2%2FIr%2FvZQPFzvWCSiq88kQ9aFWTUAg67Z0o44sREdLohJ2qbAEF8tRcAQ6wo4BMYa1GlZgJv%2Bvta0ufzNOS2J1XuG3McHoKhg46EitsEflbD%2F02Andu5iYTQPqnG5vb1GI7N035xrFobUFy0UDXVUjYVggLMdn4UhYN3i8jW54xcuAep8i1F9xN4ZisszEGK0OyhmVC6lZr%2BNJdb0giCLQM8McFvkTxTUZQEnFAsjpjdkecG0p94Nz1bRB7bJ5o5gOGXJB1bHoHpoJDZNbbABALVkoSxx5nM2oZe32ZnWAizIDcftCNAMrQbdXmQUCOLdDs8Ic%2FZEQ9E7to9H0unxhxVIhouXOuC%2F89l%2BRXQ5XTvnWLB7yEwp5JU9v1OI1q4CADsL7O6Vi0sV2nBUeiie9WoUj%2FhqDu0D7VP%2F%2Fg1TvpogXBhCx2TmE7XUAL9xbn7pr%2FLEymyD2xdGEkYAFei%2FI2%2FhLakZC%2BfyVg%2BemJNUWpLaFGNCr0X0sdhp5RcwjY1rz1GE2NrLAMX1ra4gwaOJe9YCfOy%2FftEdCNy6L3j4%2FHWb%2F%2FAUbwRcCq6SK2%2Fo5yB8SyYDXQS6vZFwwzHcZPnbLPGy0KaA7R8FqrQour6OZ4ukNWkSlHitnPReAno3XKN5GPrLDzuLKYCSrTwj91JYCjxB5Itvs5TsMlh94qtofr5%2BWu6jn61cQPkI5n5%2B9SgpGEwwvDKY4AQC1AZq8tgb4Y4D0Ue2j%2B2WS0T9iaR3AVEr9Lz3JKwxZFZSGc1tLhyBfa7r1htEnHWa30MkbsDSbLWB7X7HpVkZA%2FXyPncCVnWsQCUCb%2BrKz4MxGWDeQ9g%2BhFjE%2Bi2NJ09AkcPTgMHisSznyk%2B2U%2BPDEeC2ZNkLBckn4BNeZVgfDuw4xExmIwfNunfV2%2BEPJku6JQwqt313ZGot9tEALSO9bbuNbBB6Cw1rIy8af4WJu53ZE4lGCWi%2FnJLupxmohUu%2BbF%2FZ8snizCPembOs8uWFg4dxck17ZavBfryzDameLCIqWtmyiW6mjaXZ6yOLl7uWEmtiN2Zdd50pKYhAYVG3lrgzidtJPyGgI7IoEqELY3FEJV5CC9xUd6DCCeGFn5vF%2FQ%2BgaoPV4Uo0XTfli6l7wKP8wfSOnRSVjvVOIcj2oHi54Bdb3qb8u4jilwWCasqvEd8Fxv4xOHA0zr%2FOjZVPevaA4K6XdMIpFDp6IcCS8iGiy7EArRuSgJ0BMw4O1XNMP8g3r9bdQcpbQNBGTiLIPA9LtPYATT9TPuY9MHRrkGM4HTmU1YLjbFUBK82Y3btZSQZhUnTAJ5qan23QLdrDkb5NJsrApQ1iz3zrrU0lz%2FSeD4Yzg8yIXLNw6Kny8OLXEWUMS5I%2F%2FRx1%2BvUVX2RqsNYhV%2BeoNdLQ5Y4Jyny8fJeaaHPnwGTPOdW0jmMc%2FtUVx3%2FHno4jMekgVKBLJJLQEmk6xZzbctI5EruYyP81Wh2WWhV5UtBFU7OhLWls%2BeShpBICvxiApaujCqTZt1RGhUucZQHCbmmwQTZwzqEK7%2Fu7r04B9YelUEo5dTmPG%2BH6%2FvNnZaLYQ%2F9EfYY67oMA4i0KuKfFlGk2c0vuowyNYC2QDwjRt7VIww5cahVHznAMBtb%2Fb7FcUE5GW%2B1%2BLqXxy4ddrAM6X%2Bxc6i3G5rOcbEjcbvyFeE0Slkvv9iZCVMVSfgLQj3riVGgz6YUXRSXRwUPpdY48C8iHkzGVlV2h7fMBwxOIrEcrR97SL0TMntdwW8i6Uo53mCmu5X6&__VIEWSTATEGENERATOR=9A260B35&__AjaxControlToolkitCalendarCssLoaded=&__VIEWSTATEENCRYPTED=&__ASYNCPOST=true&ctl00%24MainContent%24ctl26%24btnSearchByInfo=%D8%A8%D8%AD%D8%AB")

		req, err := http.NewRequest("POST", "https://candidaturebac.men.gov.ma/CandidatBac/Pages/Bac/DemandesCandidature/DemandeCandidature.aspx", body)
		if err != nil {
			// Error wil making request
			c.JSON(500, gin.H{"error": "Cannot get Data please try later"})
			return
		}
		req.Header.Set("Accept", "*/*")
		req.Header.Set("Accept-Language", "en-US,en;q=0.9")
		req.Header.Set("Cache-Control", "no-cache")
		req.Header.Set("Connection", "keep-alive")
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
		req.Header.Set("Cookie", "ASP.NET_SessionId=wauvd0ktl5rke2d1akkmlcbq")
		req.Header.Set("Origin", "https://candidaturebac.men.gov.ma")
		req.Header.Set("Referer", "https://candidaturebac.men.gov.ma/CandidatBac/Pages/Bac/DemandesCandidature/DemandeCandidature.aspx")
		req.Header.Set("Sec-Fetch-Dest", "empty")
		req.Header.Set("Sec-Fetch-Mode", "cors")
		req.Header.Set("Sec-Fetch-Site", "same-origin")
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36 Edg/114.0.1823.58")
		req.Header.Set("X-Microsoftajax", "Delta=true")
		req.Header.Set("X-Requested-With", "XMLHttpRequest")
		req.Header.Set("Sec-Ch-Ua", "\"Not.A/Brand\";v=\"8\", \"Chromium\";v=\"114\", \"Microsoft Edge\";v=\"114\"")
		req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
		req.Header.Set("Sec-Ch-Ua-Platform", "\"Windows\"")

		// make a client with short timeout because in serverless time cost money
		client := &http.Client{
			Timeout: 5 * time.Second,
		}

		resp, err := client.Do(req)
		if err != nil {
			c.JSON(500, gin.H{"error": "Cannot get Data please try later"})
			return
		}
		defer resp.Body.Close()

		bytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			c.JSON(500, gin.H{"error": "Cannot turn body to bytes"})
			return
		}

		// convert the byte slice to a string
		c.String(http.StatusOK, string(bytes))
	})

}

func init() {
	app = gin.New()
	r := app.Group("/api")
	myRoute(r)

}

// ADD THIS SCRIPT
func Handler(w http.ResponseWriter, r *http.Request) {
	app.ServeHTTP(w, r)
}
