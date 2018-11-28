package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/Pallinder/go-randomdata"
	seeds "github.com/fluxynet/gocipe/util/seeds"
	"github.com/gosimple/slug"
	"github.com/icrowley/fake"
	"github.com/satori/go.uuid"
)

//URL represents a seed type Url
const URL string = "Url"

// Status represents a seed type Status
const Status string = "Status"

// ImageURL represents a seed type ImageURL
const ImageURL string = "ImageUrl"

//Timestamp represents a seed type
const Timestamp string = "Timestamp"

//CustomList represents seed type
const CustomList string = "CustomList"

//Slug represents seed type
const Slug string = "Slug"

//Time represents seed type
const Time string = "Time"

//Number represents seed type
const Number string = "Number"

//Dob represents seed type
const Dob string = "Dob"

//UUID represents seed type
const UUID string = "UUID"

//Brand represents seed type
const Brand string = "Brand"

//Character represents seed type
const Character string = "Character"

//Characters represents seed type
const Characters string = "Characters"

//CharactersN represents seed type
const CharactersN string = "CharactersN"

//City represents seed type
const City string = "City"

//Color represents seed type
const Color string = "Color"

//Company represents seed type
const Company string = "Company"

//Continent represents seed type
const Continent string = "Continent"

//Country represents seed type
const Country string = "Country"

//CreditCardNum represents seed type
const CreditCardNum string = "CreditCardNum"

//CreditCardType represents seed type
const CreditCardType string = "CreditCardType"

//Currency represents seed type
const Currency string = "Currency"

//CurrencyCode represents seed type
const CurrencyCode string = "CurrencyCode"

//Day represents seed type
const Day string = "Day"

//Digits represents seed type
const Digits string = "Digits"

//DigitsN represents seed type
const DigitsN string = "DigitsN"

//DomainName represents seed type
const DomainName string = "DomainName"

//DomainZone represents seed type
const DomainZone string = "DomainZone"

//EmailAddress represents seed type
const EmailAddress string = "EmailAddress"

//EmailBody represents seed type
const EmailBody string = "EmailBody"

//EmailSubject represents seed type
const EmailSubject string = "EmailSubject"

//FemaleFirstName represents seed type
const FemaleFirstName string = "FemaleFirstName"

//FemaleFullName represents seed type
const FemaleFullName string = "FemaleFullName"

//FemaleFullNameWithPrefix represents seed type
const FemaleFullNameWithPrefix string = "FemaleFullNameWithPrefix"

//FemaleFullNameWithSuffix represents seed type
const FemaleFullNameWithSuffix string = "FemaleFullNameWithSuffix"

//FemaleLastName represents seed type
const FemaleLastName string = "FemaleLastName"

//FemalePatronymic represents seed type
const FemalePatronymic string = "FemalePatronymic"

//FirstName represents seed type
const FirstName string = "FirstName"

//FullName represents seed type
const FullName string = "FullName"

//FullNameWithPrefix represents seed type
const FullNameWithPrefix string = "FullNameWithPrefix"

//FullNameWithSuffix represents seed type
const FullNameWithSuffix string = "FullNameWithSuffix"

//Gender represents seed type
const Gender string = "Gender"

//GenderAbbrev represents seed type
const GenderAbbrev string = "GenderAbbrev"

//GetLangs represents seed type
const GetLangs string = "GetLangs"

//HexColor represents seed type
const HexColor string = "HexColor"

//HexColorShort represents seed type
const HexColorShort string = "HexColorShort"

//IPv4 represents seed type
const IPv4 string = "IPv4"

//IPv6 represents seed type
const IPv6 string = "IPv6"

//Industry represents seed type
const Industry string = "Industry"

//JobTitle represents seed type
const JobTitle string = "JobTitle"

//Language represents seed type
const Language string = "Language"

//LastName represents seed type
const LastName string = "LastName"

//Latitude represents seed type
const Latitude string = "Latitude"

//LatitudeDegrees represents seed type
const LatitudeDegrees string = "LatitudeDegrees"

//LatitudeDirection represents seed type
const LatitudeDirection string = "LatitudeDirection"

//LatitudeMinutes represents seed type
const LatitudeMinutes string = "LatitudeMinutes"

//LatitudeSeconds represents seed type
const LatitudeSeconds string = "LatitudeSeconds"

//Longitude represents seed type
const Longitude string = "Longitude"

//LongitudeDegrees represents seed type
const LongitudeDegrees string = "LongitudeDegrees"

//LongitudeDirection represents seed type
const LongitudeDirection string = "LongitudeDirection"

//LongitudeMinutes represents seed type
const LongitudeMinutes string = "LongitudeMinutes"

//LongitudeSeconds represents seed type
const LongitudeSeconds string = "LongitudeSeconds"

//MaleFirstName represents seed type
const MaleFirstName string = "MaleFirstName"

//MaleFullName represents seed type
const MaleFullName string = "MaleFullName"

//MaleFullNameWithPrefix represents seed type
const MaleFullNameWithPrefix string = "MaleFullNameWithPrefix"

//MaleFullNameWithSuffix represents seed type
const MaleFullNameWithSuffix string = "MaleFullNameWithSuffix"

//MaleLastName represents seed type
const MaleLastName string = "MaleLastName"

//MalePatronymic represents seed type
const MalePatronymic string = "MalePatronymic"

//Model represents seed type
const Model string = "Model"

//Month represents seed type
const Month string = "Month"

//MonthNum represents seed type
const MonthNum string = "MonthNum"

//MonthShort represents seed type
const MonthShort string = "MonthShort"

//Paragraph represents seed type
const Paragraph string = "Paragraph"

//Paragraphs represents seed type
const Paragraphs string = "Paragraphs"

//ParagraphsN represents seed type
const ParagraphsN string = "ParagraphsN"

//Password represents seed type
const Password string = "Password"

//Patronymic represents seed type
const Patronymic string = "Patronymic"

//Phone represents seed type
const Phone string = "Phone"

//Product represents seed type
const Product string = "Product"

//ProductName represents seed type
const ProductName string = "ProductName"

//Sentence represents seed type
const Sentence string = "Sentence"

//Sentences represents seed type
const Sentences string = "Sentences"

//SentencesN represents seed type
const SentencesN string = "SentencesN"

//SimplePassword represents seed type
const SimplePassword string = "SimplePassword"

//State represents seed type
const State string = "State"

//StateAbbrev represents seed type
const StateAbbrev string = "StateAbbrev"

//Street represents seed type
const Street string = "Street"

//StreetAddress represents seed type
const StreetAddress string = "StreetAddress"

//Title represents seed type
const Title string = "Title"

//TopLevelDomain represents seed type
const TopLevelDomain string = "TopLevelDomain"

//UserAgent represents seed type
const UserAgent string = "UserAgent"

//UserName represents seed type
const UserName string = "UserName"

//WeekDay represents seed type
const WeekDay string = "WeekDay"

//WeekDayShort represents seed type
const WeekDayShort string = "WeekDayShort"

//WeekdayNum represents seed type
const WeekdayNum string = "WeekdayNum"

//Word represents seed type
const Word string = "Word"

//Words represents seed type
const Words string = "Words"

//WordsN represents seed type
const WordsN string = "WordsN"

//Year represents seed type
const Year string = "Year"

//Zip represents seed type
const Zip string = "Zip"

//TypeVarchar represents a seed of type Varchar
const TypeVarchar string = "varchar"

//TypeChar represents a seed of type Char
const TypeChar string = "char"

//TypeInt represents a seed of type Int
const TypeInt string = "int"

//TypeTimeStamp represents a seed of type TimeStamp
const TypeTimeStamp string = "timestamp"

//TypeBoolean represents a seed of type Boolean
const TypeBoolean string = "boolean"

//TypeText represents a seed of type Text
const TypeText string = "text"

//Seed implements the seeder
type Seed struct {
	Type    string `json:"type"`
	Options Option `json:"options"`
}

//Option represents the
type Option struct {
	Blank      string       `json:"blank"`
	CreditCard string       `json:"credit_card"`
	CustomList []string     `json:"custom_list"`
	Number     NumberOpts   `json:"number"`
	Date       DateOpts     `json:"date"`
	Datetime   DateTimeOpts `json:"datetime"`
	Time       TimeOpts     `json:"time"`
	Password   PasswordOpts `json:"password"`
}

//TimeOpts implements options for Time
type TimeOpts struct {
	From int    `json:"from"`
	To   int    `json:"to"`
	Val  string `json:"value"`
}

//NumberOpts implements option for NumberOpts
type NumberOpts struct {
	Min int `json:"min"`
	Max int `json:"max"`
	Val int `json:"val"`
}

//PasswordOpts implements option for PasswordOpts
type PasswordOpts struct {
	Min                    int  `json:"min"`
	Max                    int  `json:"max"`
	AllowUpper             bool `json:"allow_upper"`
	AllowNumeric           bool `json:"allow_numeric"`
	AllowSpecialCharacters bool `json:"allow_special_characters"`
}

//DateTimeOpts implements option for DateTimeOpts
type DateTimeOpts struct {
	Now        bool   `json:"now"`
	FromYY     int    `json:"from_yy"`
	ToYY       int    `json:"to_yy"`
	FromTimeHH int    `json:"from_time_hh"`
	ToTimeHH   int    `json:"to_time_hh"`
	Format     string `json:"format"`
	Timezone   string `json:"timezone"`
}

//DateOpts implements option for DateOpts
type DateOpts struct {
	From   int    `json:"from"`
	To     int    `json:"to"`
	Format string `json:"format"`
}

//Record implements a map[][]
type Record map[string]interface{}

//GenerataSeeds generate dummy data
func GenerataSeeds(r *Recipe) []string {
	var data = make(map[string][]Record)
	var slugs []map[string]string

	for _, entity := range r.Entities {

		var seeds Seed
		seeds.Type = "UUID"

		entity.fields["ID"].Seed = seeds

		var seed1 Seed
		seed1.Type = "Status"
		entity.fields["Status"].Seed = seed1

	}

	for _, entity := range r.Entities {
		var records []Record
		for i := 0; i < entity.SeedCount; i++ {
			var record Record = make(map[string]interface{})
			for _, field := range entity.Fields {
				record[field.schema.Field] = getData(field, entity, &slugs, record)
			}

			records = append(records, record)
		}
		data[entity.Table] = records
	}

	//Cater for relationships
	for _, items := range r.Entities {
		//Check type of relation
		for _, rel := range items.Relationships {
			// fmt.Println(rel.Type)

			switch rel.Type {
			case RelationshipTypeOneOneOwner:
				ThisEntity := getEntityIDs(data, rel, true, "")
				// lenEntity := len(ThatEntity)

				field := fmt.Sprintf("%s_id", strings.ToLower(rel.Entity))
				var i = 0
				for _, datum := range data[items.Table] {
					datum[field] = ThisEntity[i]
					i++
				}

			case RelationshipTypeManyOne:
				//get len of entity having one-many with this entity
				lenEntity := len(data[rel.related.Table])
				for _, datum := range data[items.Table] {
					randKey := fake.Year(0, lenEntity-1)
					datum[rel.related.Table] = data[rel.related.Table][randKey]["id"]
				}

			case RelationshipTypeManyManyOwner:
				//Many many owner (list of ids)
				ThisEntity := getEntityIDs(data, rel, true, "")
				//Many many inverse (list of ids)
				ThatEntity := getEntityIDs(data, rel, false, items.Table)

				var records []Record

				var max = rel.RelSeeder.MaxPerEntity - 1

				fmt.Println(max)
				//Table many many generated
				for k := range ThisEntity {
					numRel := fake.Year(0, max)
					maxNum1 := len(ThatEntity) - 1 //For entity (inverse)
					for i := 0; i < numRel; i++ {
						var record Record = make(map[string]interface{})

						randNum1 := fake.Year(0, maxNum1)

						//Get a random id of entity (own)
						record[rel.ThisID] = ThisEntity[k]
						//Get a random id of entity (inverse)
						record[rel.ThatID] = ThatEntity[randNum1]
						records = append(records, record)
					}
				}
				data[rel.JoinTable] = records
			}
		}
	}

	filename, _ := GetAbsPath("schema/seeder.json")

	d, _ := json.Marshal(data)
	_ = ioutil.WriteFile(filename, d, 0644)

	allTypes := getAllDataWithType(r, data)
	var statements []string

	for _, items := range r.Entities {
		statements = sqlization(data, allTypes, items.Table, statements)
	}

	entityNames := getEntityNames(r)

	var mnymnytbls []string
	for key := range data {
		isFalse := contains(entityNames, key)

		if isFalse == false {
			mnymnytbls = append(mnymnytbls, key)
		}
	}

	for _, manymanyT := range mnymnytbls {
		statements = sqlization(data, allTypes, manymanyT, statements)
	}

	return statements
}

func getAllDataWithType(f *Recipe, data map[string][]Record) map[string][]map[string]string {
	var datatype = make(map[string][]map[string]string)
	for _, entity := range f.Entities {
		var dataRec []map[string]string
		for _, field := range entity.Fields {
			var datatypeRec = make(map[string]string)
			datatypeRec[field.schema.Field] = field.schema.Type
			dataRec = append(dataRec, datatypeRec)
		}
		datatype[entity.Table] = dataRec

		for _, rel := range entity.Relationships {
			switch rel.Type {
			case RelationshipTypeManyOne:

				var datatypeRec = make(map[string]string)
				datatypeRec[rel.related.Table] = "VARCHAR(255)"
				dataRec = append(dataRec, datatypeRec)
				datatype[entity.Table] = dataRec

			case RelationshipTypeManyManyOwner:
				var datatypeRec = make(map[string]string)
				datatypeRec[rel.ThisID] = "VARCHAR(255)"
				dataRec = append(dataRec, datatypeRec)

				datatypeRec[rel.ThatID] = "VARCHAR(255)"
				dataRec = append(dataRec, datatypeRec)

				// str := getMnyMnyTblName(entity.Table, rel.related.Table)

				datatype[rel.JoinTable] = dataRec
			}

		}

	}

	return datatype
}

func checkType(allTypes map[string][]map[string]string, tablename string, field string, val interface{}) string {
	var fieldtype string

	for _, v := range allTypes[tablename] {
		fieldtype = v[field]

		var s = strings.ToLower(fieldtype)
		if strings.Contains(s, TypeVarchar) || strings.Contains(s, TypeChar) || strings.Contains(s, TypeTimeStamp) || strings.Contains(s, TypeText) {
			return fmt.Sprintf(" '%s'", val)
		} else if strings.Contains(s, TypeInt) {
			return strconv.Itoa(val.(int))
		}
	}

	r := fmt.Sprintf("%s", val)
	return r
}

func sqlization(data map[string][]Record, allTypes map[string][]map[string]string, entityName string, statements []string) []string {

	for _, record := range data[entityName] {
		var vals []string
		var fs []string
		for k, v := range record {
			dta := checkType(allTypes, entityName, k, v)

			// //values
			fs = append(fs, k)
			vals = append(vals, dta)
		}

		concatFields := strings.Join(fs, ", ")
		concatVals := strings.Join(vals, ",")
		stmt := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s);", entityName, concatFields, concatVals)
		stmt = fmt.Sprintf("%s\n", stmt)

		statements = append(statements, stmt)
	}

	return statements
}

func getEntityIDs(data map[string][]Record, rel Relationship, owner bool, tablename string) []string {
	var ids []string

	if owner == true {
		for _, itm := range data[rel.related.Table] { // for p in Person
			id := fmt.Sprintf("%v", itm["id"])
			ids = append(ids, id)
		}
	} else {
		for _, itm := range data[tablename] { // for p in Person
			id := fmt.Sprintf("%v", itm["id"])
			ids = append(ids, id)
		}
	}

	return ids
}

func getAllIds(data map[string][]Record, rel Relationship, owner bool, tablename string) map[string]string {
	ids := make(map[string]string)

	if owner == true {
		for _, itm := range data[rel.related.Table] { // for p in Person
			id := fmt.Sprintf("%v", itm["id"])
			ids[id] = id
		}
	} else {
		for _, itm := range data[tablename] { // for p in Person
			id := fmt.Sprintf("%v", itm["id"])
			ids[id] = id
		}
	}

	return ids
}

func getEntityNames(f *Recipe) []string {
	var names []string
	for _, entity := range f.Entities {
		names = append(names, entity.Table)
	}
	return names
}

func contains(names []string, ele string) bool {
	for _, item := range names {
		if item == ele {
			return true
		}
	}
	return false
}

func getMnyMnyTblName(tbl1, tbl2 string) string {
	var t1 = strings.ToLower(tbl1)
	var t2 = strings.ToLower(tbl2)
	if t1 < t2 {
		return fmt.Sprintf("%ss_%ss", t1, t2)
	}
	return fmt.Sprintf("%ss_%ss", t2, t1)
}

func getDummyDate(fromYY int, toYY int) string {
	if fromYY == 0 || toYY == 0 {
		now := time.Now()
		return now.Format("02-01-2006")
	}

	d := fake.Year(1, 28)
	y := fake.Year(fromYY, toYY)
	mn := fake.Year(1, 12)

	if len(strconv.Itoa(d)) == 1 && len(strconv.Itoa(mn)) == 1 {
		return fmt.Sprintf("0%d-0%d-%d", d, mn, y)
	} else if len(strconv.Itoa(d)) == 2 && len(strconv.Itoa(mn)) == 1 {
		return fmt.Sprintf("%d-0%d-%d", d, mn, y)
	} else if len(strconv.Itoa(d)) == 1 && len(strconv.Itoa(mn)) == 2 {
		return fmt.Sprintf("0%d-%d-%d", d, mn, y)
	} else {
		return fmt.Sprintf("%d-%d-%d", d, mn, y)
	}
}

func getDummyTime(from int, to int, value string) string {

	if value != "" {
		return value
	}
	if from == 0 {
		//load default value
		from = 23
	}

	if to == 0 {
		//load default value
		to = 59
	}

	hh := randomdata.Number(0, from)
	mm := randomdata.Number(0, to)

	if len(strconv.Itoa(hh)) == 1 && len(strconv.Itoa(mm)) == 1 {
		return fmt.Sprintf("0%d:0%d:00", hh, mm)
	} else if len(strconv.Itoa(hh)) == 2 && len(strconv.Itoa(mm)) == 1 {
		return fmt.Sprintf("%d:0%d:00", hh, mm)
	} else if len(strconv.Itoa(hh)) == 1 && len(strconv.Itoa(mm)) == 2 {
		return fmt.Sprintf("0%d:%d:00", hh, mm)
	} else {
		return fmt.Sprintf("%d:%d:00", hh, mm)
	}
}

//RandString generated a unique number of N length
func RandString(n int) string {
	const letterBytes string = "1234567890"
	const (
		letterIdxBits = 3                    // 3 bits to represent a number index
		letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
		letterIdxMax  = 63 / letterIdxBits   // # of number indices fitting in 63 bits
	)

	var src = rand.NewSource(time.Now().UnixNano())
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

func getData(field Field, entity Entity, slugs *[]map[string]string, record map[string]interface{}) interface{} {
	var s = make(map[string]string)

	switch field.Seed.Type {

	case Timestamp:

		dummydate := getDummyDate(field.Seed.Options.Datetime.FromYY, field.Seed.Options.Datetime.ToYY)
		dummytime := getDummyTime(field.Seed.Options.Datetime.FromTimeHH, 0, "")

		if field.Seed.Options.Datetime.Timezone != "" {
			return fmt.Sprintf("%s %s%s", dummydate, dummytime, field.Seed.Options.Datetime.Timezone)
		}

		return fmt.Sprintf("%s %s", dummydate, dummytime)

	case Number:
		return randomdata.Number(field.Seed.Options.Number.Min, field.Seed.Options.Number.Max)

	case CustomList:
		options := field.Seed.Options.CustomList
		length := len(options) - 1

		choice := randomdata.Number(0, length)
		return options[choice]

	case Slug:
		str := fmt.Sprintf("%v", record[entity.Slug])
		if len(strings.Trim(str, " ")) != 0 {
			str = slug.Make(str)

			//check if slug exits
			exist := s[str]
			if exist == "" {
				s[str] = str
				*slugs = append(*slugs, s)

			} else {
				// append a unique value to the slug
				str = fmt.Sprintf("%s-%s", str, RandString(4))
				s[str] = str
				*slugs = append(*slugs, s)
			}
			return str
		}

	case Time:
		return getDummyTime(field.Seed.Options.Time.From, field.Seed.Options.Time.To, field.Seed.Options.Time.Val)
	case Dob:
		var d int
		d = fake.Year(1, 28)
		y := fake.Year(field.Seed.Options.Date.From, field.Seed.Options.Date.To)
		mn := fake.Year(1, 12)
		var str string
		if len(strconv.Itoa(d)) == 1 && len(strconv.Itoa(mn)) == 1 {
			str = fmt.Sprintf("0%d-0%d-%d", d, mn, y)
		} else if len(strconv.Itoa(d)) == 2 && len(strconv.Itoa(mn)) == 1 {
			str = fmt.Sprintf("%d-0%d-%d", d, mn, y)
		} else if len(strconv.Itoa(d)) == 1 && len(strconv.Itoa(mn)) == 2 {
			str = fmt.Sprintf("0%d-%d-%d", d, mn, y)
		} else {
			str = fmt.Sprintf("%d-%d-%d", d, mn, y)
		}

		return str
	case UUID:
		return uuid.NewV4()
	case Brand:
		return fake.Brand()
	case Character:
		return fake.Character()
	case Characters:
		return fake.Characters()
	case CharactersN:
		return fake.CharactersN(field.Seed.Options.Number.Val)
	case City:
		return fake.City()
	case Color:
		return fake.Color()
	case Company:
		return fake.Company()
	case Continent:
		return fake.Continent()
	case Country:
		return fake.Country()
	case CreditCardNum:
		return fake.CreditCardNum(field.Seed.Options.CreditCard)
	case CreditCardType:
		return fake.CreditCardType()
	case Currency:
		return fake.Currency()
	case CurrencyCode:
		return fake.CurrencyCode()
	case Day:
		return fake.Day()
	case Digits:
		return fake.Digits()
	case DigitsN:
		return fake.DigitsN(field.Seed.Options.Number.Val)
	case DomainName:
		return fake.DomainName()
	case DomainZone:
		return fake.DomainZone()
	case EmailAddress:
		return fake.EmailAddress()
	case EmailBody:
		return fake.EmailBody()
	case EmailSubject:
		return fake.EmailSubject()
	case FemaleFirstName:
		return fake.FemaleFirstName()
	case FemaleFullName:
		return fake.FemaleFullName()
	case FemaleFullNameWithPrefix:
		return fake.FemaleFullNameWithPrefix()
	case FemaleFullNameWithSuffix:
		return fake.FemaleFullNameWithSuffix()
	case FemaleLastName:
		return fake.FemaleLastName()
	case FemalePatronymic:
		return fake.FemalePatronymic()
	case FirstName:
		return fake.FirstName()
	case FullName:
		return fake.FullName()
	case FullNameWithPrefix:
		return fake.FullNameWithPrefix()
	case FullNameWithSuffix:
		return fake.FullNameWithSuffix()
	case Gender:
		return fake.Gender()
	case GenderAbbrev:
		return fake.GenderAbbrev()
	case GetLangs:
		return fake.GetLangs()
	case HexColor:
		return fake.HexColor()
	case HexColorShort:
		return fake.HexColorShort()
	case ImageURL:
		return seeds.RandDummyImageURL()
	case IPv4:
		return fake.IPv4()
	case IPv6:
		return fake.IPv6()
	case Industry:
		return fake.Industry()
	case JobTitle:
		return fake.JobTitle()
	case Language:
		return fake.Language()
	case LastName:
		return fake.LastName()
	case Latitude:
		return fake.Latitude()
	case LatitudeDegrees:
		return fake.LatitudeDegrees()
	case LatitudeDirection:
		return fake.LatitudeDirection()
	case LatitudeMinutes:
		return fake.LatitudeMinutes()
	case LatitudeSeconds:
		return fake.LatitudeSeconds()
	case Longitude:
		return fake.Longitude()
	case LongitudeDegrees:
		return fake.LongitudeDegrees()
	case LongitudeDirection:
		return fake.LongitudeDirection()
	case LongitudeMinutes:
		return fake.LongitudeMinutes()
	case LongitudeSeconds:
		return fake.LongitudeSeconds()
	case MaleFirstName:
		return fake.MaleFirstName()
	case MaleFullName:
		return fake.MaleFullName()
	case MaleFullNameWithPrefix:
		return fake.MaleFullNameWithPrefix()
	case MaleFullNameWithSuffix:
		return fake.MaleFullNameWithSuffix()
	case MaleLastName:
		return fake.MaleLastName()
	case MalePatronymic:
		return fake.MalePatronymic()
	case Model:
		return fake.Model()
	case Month:
		return fake.Month()
	case MonthNum:
		return fake.MonthNum()
	case MonthShort:
		return fake.MonthShort()
	case Paragraph:
		return fake.Paragraph()
	case Paragraphs:
		return fake.Paragraphs()
	case ParagraphsN:
		return fake.ParagraphsN(field.Seed.Options.Number.Val)
	case Password:
		return fake.Password(field.Seed.Options.Password.Min, field.Seed.Options.Password.Max, field.Seed.Options.Password.AllowUpper, field.Seed.Options.Password.AllowNumeric, field.Seed.Options.Password.AllowSpecialCharacters)
	case Patronymic:
		return fake.Patronymic()
	case Phone:
		return fake.Phone()
	case Product:
		return fake.Product()
	case ProductName:
		return fake.ProductName()
	case Sentence:
		return fake.Sentence()
	case Sentences:
		return fake.Sentences()
	case SentencesN:
		return fake.SentencesN(field.Seed.Options.Number.Val)
	case SimplePassword:
		return fake.SimplePassword()
	case State:
		return fake.State()
	case Status:
		return "Published"
	case StateAbbrev:
		return fake.StateAbbrev()
	case Street:
		return fake.Street()
	case StreetAddress:
		return fake.StreetAddress()
	case Title:
		return fake.Title()
	case TopLevelDomain:
		return fake.TopLevelDomain()
	case URL:
		return seeds.RandDummyURL()
	case UserAgent:
		return fake.UserAgent()
	case UserName:
		return fake.UserName()
	case WeekDay:
		return fake.WeekDay()
	case WeekDayShort:
		return fake.WeekDayShort()
	case WeekdayNum:
		return fake.WeekdayNum()
	case Word:
		return fake.Word()
	case Words:
		return fake.Words()
	case WordsN:
		return fake.WordsN(field.Seed.Options.Number.Val)
	case Year:
		return fake.Year(field.Seed.Options.Number.Min, field.Seed.Options.Number.Max)
	case Zip:
		return fake.Zip()
	default:
		log.Printf("Field: %s, does not have type %s", field.schema.Field, field.Seed.Type)
		return ""
	}

	return nil
}
