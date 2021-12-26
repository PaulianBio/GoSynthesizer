package synthesizer

const KAKAOI_AUTH_KEY string = "9b8c28784daa4d246244263249e7f956"

func CreateKakaoTTSRequestor() *KakaoTTSRequestor {
	requestor := new(KakaoTTSRequestor)
	requestor.KakaoTTSRequest.RequestApi = RequestApi{
		Method: KAKAOI_METHOD,
		Scheme: KAKAOI_SCHEME,
		Host:   KAKAOI_HOST,
		Path:   KAKAOI_PATH,
	}
	requestor.ContentType = KAKAOI_CONTENT_TYPE
	requestor.Authorization = KAKAOI_AUTH_KEY
	return requestor
}
