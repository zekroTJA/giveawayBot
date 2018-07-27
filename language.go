package main

import (
	"io/ioutil"
	"github.com/go-yaml/yaml"
)

// language contains texts in set 
// language to display in imessages
type Language struct {
	Commands struct {
		Giveaway struct {
			InvalidInput,		
			EnterContent,			
			EnterWinMessage,			
			EnterParticipantsNumber, 
			EnterExpireTime,			
			EnterChannelResolvable,	
			CreatingFailed,			
			Created string
		}
	}
	Classes struct {
		Giveaway struct {
			ActiveMessage struct {
				Title,
				ParticipateInfo,
				Expires string
			}
			ClosedMessage struct {
				Title,
				Winners,
				NoParticipants,
				Expired string
			}
			CreatorDM struct {
				NoParticipations,
				Final string
			}
			Notifications struct {
				MultiParticipation,
				Participated string
			}
		}
	}
}

// NewLanguage load the entered language file set
// in config
func NewLanguage(config *Config) (*Language, error) {
	b_data, err := ioutil.ReadFile(config.Data.Language + ".yaml")
	if err != nil {
		return nil, err
	}
	data := &Language{}
	err = yaml.Unmarshal(b_data, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}