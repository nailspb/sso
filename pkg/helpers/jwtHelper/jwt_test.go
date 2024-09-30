package jwtHelper

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/ssh"
	"strconv"
	"testing"
)

type caseStruct struct {
	privateKey string
	publicKey  string
	token      string
	result     bool
}

func (c *caseStruct) GetRsaKey() *rsa.PrivateKey {
	key, err := ssh.ParseRawPrivateKey([]byte(c.privateKey))
	if err != nil {
		panic(err)
	}
	return key.(*rsa.PrivateKey)
}

func getCases() []caseStruct {
	return []caseStruct{
		{
			privateKey: "-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQDk92609HmGyYkX\nB4D8rKTJLfN+MwukU/Cu5WMhtxwhJCW0EN4xOjh6vzEIPDKegNDYFq8rf7HoadWX\ndDl0EkiFrwpopfonF2Zu0Rgx/Q8H7Gpwpa6fIZu/Xc8IWnPAVWmUXXRQviuuXwYd\n0bOimg2FCKByhV1pcOIox5piR8kQ36Dw30dKGFJt2kAxRX011MWz3dkAl5UShXoJ\nSP+OwR8R62pro7XdT+HfoCELALT1dU4OMZiGOt2bsW5RxE9HnpJMuHbGRJXK/3U1\niOcmW3eqBvy4IUX8PzU+cXxMeNGEZAoJYpDd6k1NWBBPZ+NY6Ye7VnE6wV7kpbXm\n73AuOOSpAgMBAAECggEASiCnSVNq8s3VocyOoH8+XFmRntZk2v9ICT9/kB9JrpsW\nz8y0OsjIF0lF1Q40LyeyNXtmm9UcIov9GCsLHL58lzc7zfSKX9SOF6t/Q1PT5XNP\nZIfnWGKXm2WCDQjHt8mPRHbbHPbsimf+QXIuT6LDZMINu/Xcb7IZri0tGMt314uK\nOXjY7K56+L8EQwjKe8n5zYatWBupZEKgM25xR3wvtW9jzzGXa8CjYzGlOYsI1pA6\nptDLyAJK1OqeYewv6zCwbDBZ33F7bSXYLcM1AtSf9ZE7OsCTXMenybnOTr1OFrhw\nMQsr08XCiJhcIqTpKLiPUmo4lYCva0Q8m9pbEi20JwKBgQDzRUSsVWRXR9YMMuGq\nkQrEAfXg702XIj0E/+O3ELenCqVXgyoz1956zkagA7LNBI7G7PYczScqMRIX7xt0\n6PAyQnWjMTsNPnGxuxTFfAb7uD+elfmtdpRUO1fVutUFWQGNibt62hl9UuvMxYQ2\nz+EAjV7h+Zy+m4nUWamb7OZ1IwKBgQDw8o3mzKDr9hkvV4Ld/xNn8zwui8EACG0I\nxCG//XL/XaS8ppeUbdhbHnC4Ez1xrgIsYnJz3VhRv+uOxedS+N795QeT+zc9N2kD\nnge5wbelqJx4orDgNICpMUSLDANPaWnGY1e1Y5EHyfstUeeiQom0+Xm6duK5qPGC\nLEKsxsHZwwKBgQDYrZjWOpTOHODtKqDPsLK7FNfxSpR8ifV5r4Ye91ftA9FzWhPL\n63lxPquvOLwYWffl/QfVbXF15hEsmj+FaTkQOxvWiDIFwNm5qV197NO3f0vDL+gc\ndeL2B4lbiDbWtYlpjQUdDofnlWULld4GiC/rsv+RHShcqeMg7d/hTyeRqwKBgFzJ\nq6fT6ay0yyIWG0mOb1S6sNRj8WEn3YVgsnaTDfQVhdk4dmssmgMNB+97SVA76I5b\nIyRHezmQJRCIWfrz6DvyNSbhuXYTnpdMBkGcvjJHampyjJbq4RlG5dR+PdAZEija\nHO63dyR+vgHH5uHvqcRNxnjuS4Wf79FnZg3PRNutAoGAIMDBFiPQYU/kyEuTMBKN\n3NV0z7Jd5ABIsqyhpt1g7v7noYX5cC7qrW+LLxvN4yYSwjzKRY5ETgZuY+d9Vq6M\n5zFSyY7wIYh9Xx/oVffbkuWBA2jzOJ0kTBTjB14Jb7dod4On+++4xCA22xFZacEt\nVIWB5ZvWDLpEr/KXEexkkkQ=\n-----END PRIVATE KEY-----",
			publicKey:  "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA5PdutPR5hsmJFweA/Kyk\nyS3zfjMLpFPwruVjIbccISQltBDeMTo4er8xCDwynoDQ2BavK3+x6GnVl3Q5dBJI\nha8KaKX6JxdmbtEYMf0PB+xqcKWunyGbv13PCFpzwFVplF10UL4rrl8GHdGzopoN\nhQigcoVdaXDiKMeaYkfJEN+g8N9HShhSbdpAMUV9NdTFs93ZAJeVEoV6CUj/jsEf\nEetqa6O13U/h36AhCwC09XVODjGYhjrdm7FuUcRPR56STLh2xkSVyv91NYjnJlt3\nqgb8uCFF/D81PnF8THjRhGQKCWKQ3epNTVgQT2fjWOmHu1ZxOsFe5KW15u9wLjjk\nqQIDAQAB\n-----END PUBLIC KEY-----",
			token:      "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJFeHBpcmVkQXQiOiIyMDI0LTA1LTIxVDE5OjMwOjEzKzAzOjAwIiwiaWQiOiJmZTYyYjhkZC1hMWZjLTQyYTMtOGM0Ni1mM2E3NDE3OTVlYTMifQ.MgBuqF_7jfhG-4XoQFJreTi09QSm3mSvZGWW9910UPsLhNcd7L2WY2dvd_GHygGPMobc14HPhhrh4jYUesfSn-puWxdIihdr_uGvX7Fiay35WSsAy_Ed3QT0710izqd9gP9-75Tr2H8usgC3LMEmaKpGIiwY1XdhVncrNFpm7zy4ck3yZGguGlMcWw0yWT7y4CvshODSme8TwqatqzjT4CXYmDEurbACA-nF08ckuX0R_FexXjXlyYqvFHOa5TlOxqU_qMiXXvaMM1Ndb9oDbtdzyrnxCS9gm-vSJ2GbOSFXrQBmLMOto3mVrJb8utvKOC_QpyD1FabcpQ3FVVmH_A",
			result:     false,
		},
		{
			privateKey: "-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQDCIPR5c6V74x8Z\nhDeqyVHKuWsptloL8LcybJ+wlHVdkiSJl25mGhXAt78GGmQkeUh6d3NTxDL5Vwao\nKMhbCixCnrJ3yU6YvhtChqE4rtzh74uZ5syjbJZO0td6WwJCs+z9EVHoQwNwrCxs\nanQWTEkjZ909zvNbWr70R4LR76Z4JMdQiwufv4FOAl3hB3vH3uXnFP8uq/DXA1nD\nik2pS8pJumQ7D1JTSXzy6Pbx9sUbEsMUuYl/G8E/8gW6MRFRDNEqwlytiHeW4T3m\ntTLB4Gx0uxdvZZ67wirNzXwGji0d1LR7AXqCa10AYl+DR83Kifl23/YUCoP7zRtJ\n2DJFTcTtAgMBAAECggEADJsmvP9HQvr98UJ+eeruUoeFX7VNdiKIxu+6U/VkBLJB\niKTKgLtXZ8cW0r7ie2Lb08YPeWRPehUOY9uYHR9VPfM/vBsYrT8FFqYW8I8cXViY\nnGMaLU+XGCHFZ9bil/HGpS1bgAVgAxBPJzAnfy4vfqrXOVUHxQooxcSiCF9zVDBU\nPQ6ruYYcGvquAqr0+TnESf4T5Yl+xwDUC7gxToUzpUjL/rgk9gbv2kPGYbO1LtTU\nMCUbzm42c2DfnuX1AFnYBbCxVbT3eje1omoraG6ehVsuFrD4CHe2Ll1tPoEksDvn\n1C+9QhpizWv2tys66mCKhgb/lzpnIP3zQZGLWzPdIQKBgQDo5rQWaSosx+Xr8UXy\nY/e8CC7FpknQ3Jt4mRap3lzdB0wZJo+EPaqyBo30VpMPN788R61Ox0CQqqp6p5bp\n27OctgEd3fnpzvCsG+fY8Uy5b34tKDFPL22tfcvccwUMWphgmiD3QhWm6bliPtew\n1NfuvVs5aFBZzta/mIinBKp/zQKBgQDVYdRiZVkQLwVr+baj7BO05p8/8oBVVbKx\nanDuanYB6H05h3mdBMSujG6zqAgO0rz1zTICSp2Vk3hrGXwwQJ4p+OKRsaCPEDA4\novodsyj5W13DM2zfOsbVSHNWdp+ZUYYprMeVDSKZ5ycQTE5r/ttEsaMT2auOUrv2\nwNPP19n5oQKBgQC1+Pwgvo6raozlKyEh4KYCDsqst97ZCGWZdtPUad9Y+4ij5hMg\nmiYy7xgfHz3MGByddsucz+ZMomyNZUBu/LB782Ev+u53lq9JaoxeXzs5cDnAArV0\nT3R8p9uPJXd+TZLdd8/mIYZVYizs7HkOu170NJOAaVbOtOPp076B8Q1eyQKBgFlX\njMVmCdRs26hJ9d7MibPEnAj6UCFqsFb4ajBpAt/pqATeZF0KEg/DXNZ8FGOgeN2x\n/K5Y74IhLNoq4YSSiaapPrQh20gLTyVnl7G3wgAl8Sw56+vLgFTs8N3S7SAUskfg\nv+/4f/RQhFqemHc/Ti+E8PLuwJXmriyyr/zmM64BAoGAbbJiSTWGyfupazMt5ccG\nhjt+olydKnOzlXI9H/TO3Q0lIeX4+B1vHeOXD+Z8Z8p/WLRbc6zuokQw1A5zJKC1\nXY5T92EOn7HDFhzu+qx/ZIVBrfqCCg/pXD+25L5CvppIYvhsqpdF3E8AnSeRLDp5\nvFvN0LM+80oqYUinH+RvPiQ=\n-----END PRIVATE KEY-----",
			publicKey:  "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAwiD0eXOle+MfGYQ3qslR\nyrlrKbZaC/C3MmyfsJR1XZIkiZduZhoVwLe/BhpkJHlIendzU8Qy+VcGqCjIWwos\nQp6yd8lOmL4bQoahOK7c4e+LmebMo2yWTtLXelsCQrPs/RFR6EMDcKwsbGp0FkxJ\nI2fdPc7zW1q+9EeC0e+meCTHUIsLn7+BTgJd4Qd7x97l5xT/Lqvw1wNZw4pNqUvK\nSbpkOw9SU0l88uj28fbFGxLDFLmJfxvBP/IFujERUQzRKsJcrYh3luE95rUyweBs\ndLsXb2Weu8Iqzc18Bo4tHdS0ewF6gmtdAGJfg0fNyon5dt/2FAqD+80bSdgyRU3E\n7QIDAQAB\n-----END PUBLIC KEY-----",
			token:      "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJFeHBpcmVkQXQiOiIyMDI0LTA1LTIxVDE5OjMwOjEzKzAzOjAwIiwiaWQiOiJmZTYyYjhkZC1hMWZjLTQyYTMtOGM0Ni1mM2E3NDE3OTVlYTMifQ.XUyQCip9DGluqWQDwRgnOqcPxbvBsazZHOjEAo2x0KZS3mrI6x2imJjsbVg-ShsZXRtMfypUc6zEFJ9R9tgZlThxNRxGW5FoSlXRIA_HMftClJLBYOfFoP99mkAr1DdboeX9h6JGhIvOo9MPE4OEgCc7qP6C_nnuCmXOlfCsjaAVSH2ADbecYgEFuXMXXvCeUbp8Pt3P1u9jCBPW3ZGVzd53epm2OuUaYL7k7w2e8ltNX0LUda3qfbBGyTps7kUf3IbcjHFJJ_9AkNgROt8RlRBDLOTC0nnaltpSoL4xcF3bK4vxORFrS6i5CSXEPEuiFqqIBDlFctd7tahs6-EKoA",
			result:     true,
		},
		{
			privateKey: "-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQDk92609HmGyYkX\nB4D8rKTJLfN+MwukU/Cu5WMhtxwhJCW0EN4xOjh6vzEIPDKegNDYFq8rf7HoadWX\ndDl0EkiFrwpopfonF2Zu0Rgx/Q8H7Gpwpa6fIZu/Xc8IWnPAVWmUXXRQviuuXwYd\n0bOimg2FCKByhV1pcOIox5piR8kQ36Dw30dKGFJt2kAxRX011MWz3dkAl5UShXoJ\nSP+OwR8R62pro7XdT+HfoCELALT1dU4OMZiGOt2bsW5RxE9HnpJMuHbGRJXK/3U1\niOcmW3eqBvy4IUX8PzU+cXxMeNGEZAoJYpDd6k1NWBBPZ+NY6Ye7VnE6wV7kpbXm\n73AuOOSpAgMBAAECggEASiCnSVNq8s3VocyOoH8+XFmRntZk2v9ICT9/kB9JrpsW\nz8y0OsjIF0lF1Q40LyeyNXtmm9UcIov9GCsLHL58lzc7zfSKX9SOF6t/Q1PT5XNP\nZIfnWGKXm2WCDQjHt8mPRHbbHPbsimf+QXIuT6LDZMINu/Xcb7IZri0tGMt314uK\nOXjY7K56+L8EQwjKe8n5zYatWBupZEKgM25xR3wvtW9jzzGXa8CjYzGlOYsI1pA6\nptDLyAJK1OqeYewv6zCwbDBZ33F7bSXYLcM1AtSf9ZE7OsCTXMenybnOTr1OFrhw\nMQsr08XCiJhcIqTpKLiPUmo4lYCva0Q8m9pbEi20JwKBgQDzRUSsVWRXR9YMMuGq\nkQrEAfXg702XIj0E/+O3ELenCqVXgyoz1956zkagA7LNBI7G7PYczScqMRIX7xt0\n6PAyQnWjMTsNPnGxuxTFfAb7uD+elfmtdpRUO1fVutUFWQGNibt62hl9UuvMxYQ2\nz+EAjV7h+Zy+m4nUWamb7OZ1IwKBgQDw8o3mzKDr9hkvV4Ld/xNn8zwui8EACG0I\nxCG//XL/XaS8ppeUbdhbHnC4Ez1xrgIsYnJz3VhRv+uOxedS+N795QeT+zc9N2kD\nnge5wbelqJx4orDgNICpMUSLDANPaWnGY1e1Y5EHyfstUeeiQom0+Xm6duK5qPGC\nLEKsxsHZwwKBgQDYrZjWOpTOHODtKqDPsLK7FNfxSpR8ifV5r4Ye91ftA9FzWhPL\n63lxPquvOLwYWffl/QfVbXF15hEsmj+FaTkQOxvWiDIFwNm5qV197NO3f0vDL+gc\ndeL2B4lbiDbWtYlpjQUdDofnlWULld4GiC/rsv+RHShcqeMg7d/hTyeRqwKBgFzJ\nq6fT6ay0yyIWG0mOb1S6sNRj8WEn3YVgsnaTDfQVhdk4dmssmgMNB+97SVA76I5b\nIyRHezmQJRCIWfrz6DvyNSbhuXYTnpdMBkGcvjJHampyjJbq4RlG5dR+PdAZEija\nHO63dyR+vgHH5uHvqcRNxnjuS4Wf79FnZg3PRNutAoGAIMDBFiPQYU/kyEuTMBKN\n3NV0z7Jd5ABIsqyhpt1g7v7noYX5cC7qrW+LLxvN4yYSwjzKRY5ETgZuY+d9Vq6M\n5zFSyY7wIYh9Xx/oVffbkuWBA2jzOJ0kTBTjB14Jb7dod4On+++4xCA22xFZacEt\nVIWB5ZvWDLpEr/KXEexkkkQ=\n-----END PRIVATE KEY-----",
			publicKey:  "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA5PdutPR5hsmJFweA/Kyk\nyS3zfjMLpFPwruVjIbccISQltBDeMTo4er8xCDwynoDQ2BavK3+x6GnVl3Q5dBJI\nha8KaKX6JxdmbtEYMf0PB+xqcKWunyGbv13PCFpzwFVplF10UL4rrl8GHdGzopoN\nhQigcoVdaXDiKMeaYkfJEN+g8N9HShhSbdpAMUV9NdTFs93ZAJeVEoV6CUj/jsEf\nEetqa6O13U/h36AhCwC09XVODjGYhjrdm7FuUcRPR56STLh2xkSVyv91NYjnJlt3\nqgb8uCFF/D81PnF8THjRhGQKCWKQ3epNTVgQT2fjWOmHu1ZxOsFe5KW15u9wLjjk\nqQIDAQAB\n-----END PUBLIC KEY-----",
			token:      "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJFeHBpcmVkQXQiOiIyMDI0LTA1LTIxVDE5OjMwOjEzKzAzOjAwIiwiaWQiOiJmZTYyYjhkZC1hMWZjLTQyYTMtOGM0Ni1mM2E3NDE3OTVlYTMifQ.2qH3o7SXEtCCSmcAD7TMfZ0jRgFaYOEjxhADrCH3-vn9i22MrIfBWpNKttScSlnhhjCu4hf3AfKofQbQIrUJRbELaOOFeHmPWsq8FoRhdPp6y1-sQ3R-DzIABsY9YTjClQkasyOobkzB1IdXuVawyY_gKcArgDdzvJNmZSwo6KNysIkRBU_XX59--zF-yDzY_HTQODfxziOie2ZLV4HmDEGWv3cWPnD5ofZkGtP99Fxj4yebSXSnwge9cZnjSP70Kg3o0Rsj6a1oMCwVk_kFR7pzsTXZbPK2zJV9MRPla-nFPNuNd_O7ET576x7ZHWlDF8gkeDdGJQcyLJaLUxTasQ",
			result:     true,
		},
		{
			privateKey: "-----BEGIN PRIVATE KEY-----\nMIIEvAIBADANBgkqhkiG9w0BAQEFAASCBKYwggSiAgEAAoIBAQDm9VXw9lNg+mkp\nwFcVUFXFcOluZMj1V+XARdJZO5wpcU0ifNaDka+4fAfwBL9Qziwn4TJoSyaiMSYx\nFj0+46XreCbYn2iz2ZaIjFz3prED3b10f/ntomaK2+njdzy+CTRmv87lAiapVZrM\nigd6j024NIuS8GbvU9GP2GOKXXAyRoJnMO8IFc5Is0yN4KElzAfTq8s0XYe+EAgq\nlMmpEoYBWZDDScecwflr4uAHC5yQqIVtSQMHxo2acSW/fUhG2HIivWOnDXSh4qn2\n7d+xYlGmoJZQ3xMJ9OMHppqSPF7A+7g1bOJfwHUVYn0/SxIypxv59NAwpA31bubQ\nJz9CX/spAgMBAAECgf9vJw/uLLkRGDEVtMr3s7n9Y7kO1vkXw6H0EXbjhbUL+FE5\noZkyR3HJGyT7bp9IRxPBw+8PuVAd684uZWs5X7fuwEudV0pWm2GVSNC8mXouT65i\nUKajmlrsXvou5e9AyFThUz8fxfu4US4sX1lu3MBzBnKUxop0fzbtQwpQgo1FT5rN\nOllVKDw+6NXp4hQbNWN6xX+jOBs3u4syameCMx2ZXgWFVLQfFFUAV8jsW158+S4K\nX7CgLj+9+PzHpz/ws8c+3IDHgHCxOmb6XmAfGJkQFZ7HxONsQV6WLR1eJm9j2rWK\nlrKoH7k1oFamftmAidUQ3BjhXLNIWnkrtLZfBbECgYEA+aVImyLJJXD8+Pa6sCLU\nGZjo/9Yqvj6IQZA5U/9MngIh/An64F7bHX89NvCCWQo3yN4SDVqU0Lalt0ZyjHzG\nVC/T9VAl3ofwXUDQxTcEIa45gFGzrknuhvNLtZq2njxXan8Nar9euPRVruQUbD4R\ni802ct21F623yQhcZBAqYQUCgYEA7NZIpNu46xLwjPVtX+3Z0BiKZV52qo73fP9g\nYSZHk3XOW4xy9Z4wdhftsWyyK3hcu93V3L5oMV2EugCyVSDnPS6GY7YyJVibIn23\nq3a7X2MYfMwxamnOsKfnKDuwqRNChpORB9XlpZDTRD0IpakVaiPD+xJAkYFPsI8l\nxfqz2tUCgYAV6SPOOeddmeUaFM2d1/C1rm4Exk9KE0LyPi5J6QZYd+dzr4yNVMX2\neMxunf1Sw0rSHmuHMIQPLXit9Ujoe6sMiIYZ6cbpGRVHmgC4znNWYWw6jvEuQt7k\niUYD0mhkyvcBKdWLoPA3W3qJtwrz1R7FHmXA/yR9x4lx44H4ZlLR4QKBgQCS2K6E\nIYR/pOenzsj5UOXbpEuzXKXhTPHg/AsLUYvRv5qqouPorSPJJT8I4qd6Uo/VIE/p\nJdo+uYiBN8tbAyK9iapkCuT+yPivoxmN4/l7xFq7jnQZUe+JEyI9jP0VaE64WKj1\nHcfdJ3YG+nzxEmynufNbKk8EqRP7GlbGcZKw3QKBgQDH8HVKhVfpboSbypWfQ4xe\nzrmGn4iOIcBkQ/wcarqhvlDyIoezX2X0+9hEg32tHjIiDlG140VsE09BP+BHdRZa\nT1u5cZZYF+obdW5BYq7GXp8v6AE6QnSI1Vl8zLJc5yWXqhkq849mRes059Cjv4jZ\nnj/G+/sveU/xWEoxFwe9/w==\n-----END PRIVATE KEY-----",
			publicKey:  "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA5vVV8PZTYPppKcBXFVBV\nxXDpbmTI9VflwEXSWTucKXFNInzWg5GvuHwH8AS/UM4sJ+EyaEsmojEmMRY9PuOl\n63gm2J9os9mWiIxc96axA929dH/57aJmitvp43c8vgk0Zr/O5QImqVWazIoHeo9N\nuDSLkvBm71PRj9hjil1wMkaCZzDvCBXOSLNMjeChJcwH06vLNF2HvhAIKpTJqRKG\nAVmQw0nHnMH5a+LgBwuckKiFbUkDB8aNmnElv31IRthyIr1jpw10oeKp9u3fsWJR\npqCWUN8TCfTjB6aakjxewPu4NWziX8B1FWJ9P0sSMqcb+fTQMKQN9W7m0Cc/Ql/7\nKQIDAQAB\n-----END PUBLIC KEY-----",
			token:      "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJFeHBpcmVkQXQiOiIyMDI0LTA1LTIxVDE5OjMwOjEzKzAzOjAwIiwiaWQiOiJmZTYyYjhkZC1hMWZjLTQyYTMtOGM0Ni1mM2E3NDE3OTVlYTMifQ.T5R0MjhJ9i1xSZpSpfIM3BjH4d23wcQl5UkpiTEqyjte9qab2HY4xNa1ADnAhp7icHaiV4k5yg6ARHf_vEO4neTa-buS2Aqs806hDci6Bc4paSEb_EOQK6RkiCQJUAlDuNLnAwP0qdNehYIKQcWdpPmFSaPzRf0lRaJWR2z01GtTLeuF8NmPDMvGifiUeloBimkXvHicymvs66etFQwkEIHmFrAWIqcqtc1FAi4M_xbUi2jLNc81dujS7RrWGCgh2y78VzvkYiL2WzwyesnwBm5kaL9US2qPab-e0hhRofdT_XUGFkg5iNdEC-1i9BMKLspcO8qigobgxO8qkv3QEw",
			result:     false,
		},
		{
			privateKey: "-----BEGIN PRIVATE KEY-----\nMIIEvwIBADANBgkqhkiG9w0BAQEFAASCBKkwggSlAgEAAoIBAQChyU81MX2Y7wtS\nG9BUw+JCZ/WdktLxfyZeExkLUn0Q4FBal67LDUmMpd5LtAc4NnRYPqAMSfSbqs5M\nKvsacBBy9Z0B+5jXJ5nbATwOy/VutdenKBF5N04dFusEswTzPZz31ARME2+0GrYd\nkjMQ2QExRYkfwjcsciwOnSqttO9bSNqxDh9rJqi/e1c8sA65g8WdLyoDEXOqnjIg\nbcjlvBHoLHTWpkGUcD93EsvHAzREvTUiTM6sD7F6WRmIIZNDIaeEPj1ce43UOKeZ\nbr+i7bN/meBX92PnxPRagpVZDJcvcNUtBlYj6GtzWNzYzJyAl3VNZKIwPWIcwixA\nWjlptHatAgMBAAECggEACXzP9Y9drXCFQi8p4DKti9XjYDZ6QtnBQv6NNy4C8hCp\nfaqMAFHa6JHYj6lyy+/bkko34gibPN6/ItXZbGM6f7daGgHqZGGo4uA3aKm6zp+A\n0kdmAO0Gmmub5CZoaahJEnp/NwPjxtTTLbhDYN9M89n/UNq3mBZo8YFzfYTRDZIH\n5FgunjsmB5EC/D2pxRuP0yA1P6sF/Qk6nKy5nECN1nMki7nT2fdMYY4Jqb+36Q/l\nOj2AYmQoWUs5UKpc9VvcZiXREF3s63p/DiHiOROxM2fynWsC039VVIwr8nVa6sYn\n2CBRL91zaas1GkgDGYdR4lL2e2PJfiAFIkBorFH/QQKBgQDkPsv8HNsgaAvT1XfL\nnlYaulGzbkHJhDMWgI+cGXIkhVLemfdcqK2GIVhmqvBpA1p8/iRF5NNoJtM+yDl/\nkBMzof2wYnyECvW47hKJKtMhpRmiy4Bbfu6MlViPmM/PGbGTjsT6DBJk00MAiB37\ni+H47NcNcZZCR45b/eRl6STE7QKBgQC1darNRI1cCAiMB3DhUiShMtQc4dPvFnen\nciaYPe5/qX3QmK4wjdSE45jGipq8JkWCD7GX46icsTjYgke4mFdS5PxHFxKE8Deo\nAQIqUL4lZ8nCmFRyItpQQ8HC1Qy2wfaRe1c8rJnXdRBUnTWwY97F/RSt+7MsB0JZ\nwip+tWwAwQKBgQCvRD11gl8N+neimhhchmCOM6bLRw5DhT2JuP0OHEgXHT3ua4KU\nZ36gMfjlFbx4lzekJa+K8FUadD3gxvXIK8Vi77CUAnylFJluAXrAU98+xb3y0Zvm\nJold8MzJr8lbudovegFuFVkGjWe0/9EuOVMzyAK3cxK0IKiDoWoCi85NXQKBgQCz\nUnYZZj6ADVxM7WmK8f9K4g0mAbHMG5rhefhUCRfxRxETnF6/ktnK/ZRT2FNzzipw\njovFe3B8cNKpe43fCYV6YNpCcrWVdEK8H0sBgEt4cam8SYdiR7kRCvSnUp2+2c2O\nFaKvTi618nTR+Y4+I2PaqvDNwuhcgUv7odsW3ri9wQKBgQDK/CH3PFWsq1JfPQUY\nzBrp2ZpxkOFVgPK0pZ0OQkwJmF8PjuBPgSjCwl+UH2sEN9a56TM1s3BufJuNK5n6\nRJlPzeGZE2q2e1fkMeu7azg2rxSgjHxEMMxal5InWoLubQOMRnOKJCvKB0vSmbaF\nm6KF0Vq/AtSFNped8RjhPtmviQ==\n-----END PRIVATE KEY-----",
			publicKey:  "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAoclPNTF9mO8LUhvQVMPi\nQmf1nZLS8X8mXhMZC1J9EOBQWpeuyw1JjKXeS7QHODZ0WD6gDEn0m6rOTCr7GnAQ\ncvWdAfuY1yeZ2wE8Dsv1brXXpygReTdOHRbrBLME8z2c99QETBNvtBq2HZIzENkB\nMUWJH8I3LHIsDp0qrbTvW0jasQ4fayaov3tXPLAOuYPFnS8qAxFzqp4yIG3I5bwR\n6Cx01qZBlHA/dxLLxwM0RL01IkzOrA+xelkZiCGTQyGnhD49XHuN1DinmW6/ou2z\nf5ngV/dj58T0WoKVWQyXL3DVLQZWI+hrc1jc2MycgJd1TWSiMD1iHMIsQFo5abR2\nrQIDAQAB\n-----END PUBLIC KEY-----",
			token:      "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJFeHBpcmVkQXQiOiIyMDI0LTA1LTIxVDE5OjMwOjEzKzAzOjAwIiwiaWQiOiJmZTYyYjhkZC1hMWZjLTQyYTMtOGM0Ni1mM2E3NDE3OTVlYTMifQ.g6Hl0PHTO61LgofQyl-Bh2taaEGsCNnBWAIYJ_eNyW1pvt_c3hzSK_YdqfEYbUJeWo299ixbKUIWp6o1lhM9Swd-QHsG-BOOYao6EjxnVnHkLqIhyR8w1dew7lvkLoitx6U4GSCtfzMOWmLjugQ8Cq2oYhkEj9kwiN3ByFNm0bZZ67u4N9hgdss6cItW_fGNrG9RUcV-98NrDcGijyNmkUXPDInRrXKkaLYM1xUu0l8bUGkdCPiThyKm2LNqqGFfl8YrAhWbItsey6JGZPc-CiI779ukFWrEqj1BDPTnlT1Tc7cWcSwVo7zXDboa---K2O39Ml_shNvVlB8Rd99ERw",
			result:     true,
		},
	}
}

func TestCreate(t *testing.T) {
	for i, c := range getCases() {
		t.Run("test "+strconv.Itoa(i), func(t *testing.T) {
			token, err := Create(c.GetRsaKey(), map[string]any{
				"ExpiredAt": "2024-05-21T19:30:13+03:00",
				"id":        "fe62b8dd-a1fc-42a3-8c46-f3a741795ea3",
			})
			if c.result {
				require.Equal(t, token, c.token)
			} else {
				require.NotEqual(t, token, c.token)
			}
			require.NoError(t, err)
		})
	}
}

func TestCheck(t *testing.T) {
	for i, c := range getCases() {
		t.Run("test "+strconv.Itoa(i), func(t *testing.T) {
			pemBlock, _ := pem.Decode([]byte(c.publicKey))
			publicKey, err := x509.ParsePKIXPublicKey(pemBlock.Bytes)
			require.NoError(t, err)
			err = Check(publicKey.(*rsa.PublicKey), c.token)
			if c.result {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func TestGetClaim(t *testing.T) {
	for i, c := range getCases() {
		t.Run("test "+strconv.Itoa(i), func(t *testing.T) {
			pemBlock, _ := pem.Decode([]byte(c.publicKey))
			publicKey, err := x509.ParsePKIXPublicKey(pemBlock.Bytes)
			assert.NoError(t, err)
			key := publicKey.(*rsa.PublicKey)
			claim, err := GetClaim(key, c.token)

			if c.result {
				assert.NoError(t, err)
				assert.Equal(t, (*claim)["ExpiredAt"], "2024-05-21T19:30:13+03:00")
				assert.Equal(t, (*claim)["id"], "fe62b8dd-a1fc-42a3-8c46-f3a741795ea3")
			} else {
				assert.Error(t, err)
			}

		})
	}
}
