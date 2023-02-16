package entities

import (
	"time"
)

type Click struct {
	Data struct {
		Id        int       `json:"id"`
		UserId    int       `json:"user_id"`
		OrderId   int       `json:"order_id"`
		Type      string    `json:"type"`
		Check     string    `json:"check"`
		CreatedAt time.Time `json:"created_at"`
	} `json:"data"`
}

type Refound struct {
	Data struct {
		Id            int    `json:"id"`
		DateTime      string `json:"date_time"`
		Amount        string `json:"amount"`
		Reason        string `json:"reason"`
		ReasonComment string `json:"reason_comment"`
		Status        string `json:"status"`
		StatusName    string `json:"status_name"`
		PaymentType   string `json:"payment_type"`
		Order         struct {
			Id       int         `json:"id"`
			OrderId  int         `json:"order_id"`
			SellerId int         `json:"seller_id"`
			AgentId  interface{} `json:"agent_id"`
			Seller   struct {
				Id        int         `json:"id"`
				ParentId  int         `json:"parent_id"`
				Name      string      `json:"name"`
				Type      string      `json:"type"`
				CompanyId int         `json:"company_id"`
				ContactId int         `json:"contact_id"`
				DeletedAt interface{} `json:"deleted_at"`
				CreatedAt interface{} `json:"created_at"`
				UpdatedAt time.Time   `json:"updated_at"`
				Contact   struct {
					Id           int         `json:"id"`
					DeletedAt    interface{} `json:"deleted_at"`
					CreatedAt    time.Time   `json:"created_at"`
					UpdatedAt    time.Time   `json:"updated_at"`
					Phone        interface{} `json:"phone"`
					Email        interface{} `json:"email"`
					Gender       interface{} `json:"gender"`
					DateBirth    interface{} `json:"date_birth"`
					PassNum      interface{} `json:"pass_num"`
					PassDate     interface{} `json:"pass_date"`
					PassIssued   interface{} `json:"pass_issued"`
					PassDivision interface{} `json:"pass_division"`
					Inn          interface{} `json:"inn"`
					Snils        interface{} `json:"snils"`
					Address      interface{} `json:"address"`
					Type         string      `json:"type"`
					Name         string      `json:"name"`
					LastName     interface{} `json:"last_name"`
					MiddleName   interface{} `json:"middle_name"`
				} `json:"contact"`
			} `json:"seller"`
			SellerName  string      `json:"seller_name"`
			SellerType  string      `json:"seller_type"`
			ContactId   int         `json:"contact_id"`
			ContactName string      `json:"contact_name"`
			ClientId    interface{} `json:"client_id"`
			ClientName  string      `json:"client_name"`
			Type        string      `json:"type"`
			PaymentType interface{} `json:"payment_type"`
			User        struct {
				Id              int         `json:"id"`
				Name            string      `json:"name"`
				Email           string      `json:"email"`
				EmailVerifiedAt time.Time   `json:"email_verified_at"`
				CreatedAt       time.Time   `json:"created_at"`
				UpdatedAt       time.Time   `json:"updated_at"`
				Position        string      `json:"position"`
				Organization    string      `json:"organization"`
				Avatar          interface{} `json:"avatar"`
				ContactId       int         `json:"contact_id"`
				SellerId        int         `json:"seller_id"`
				Timezone        string      `json:"timezone"`
				CashboxId       int         `json:"cashbox_id"`
			} `json:"user"`
			UserId  int `json:"user_id"`
			EventId int `json:"event_id"`
			Event   struct {
				Id            int    `json:"id"`
				Name          string `json:"name"`
				Datetime      string `json:"datetime"`
				Type          string `json:"type"`
				TimeLength    string `json:"timeLength"`
				AgeLimit      string `json:"ageLimit"`
				ShowId        int    `json:"show_id"`
				HallId        int    `json:"hall_id"`
				PlaceHallName string `json:"place_hall_name"`
				Promoter      struct {
					Name string `json:"name"`
				} `json:"promoter"`
				TicketFund     string `json:"ticketFund"`
				TicketProvider struct {
					Name string `json:"name"`
				} `json:"ticketProvider"`
				TicketBook int `json:"ticketBook"`
				Tickets    struct {
					Total int `json:"total"`
				} `json:"tickets"`
				Seats struct {
					Total int `json:"total"`
				} `json:"seats"`
				Place struct {
					Id       int      `json:"id"`
					Name     string   `json:"name"`
					Rows     int      `json:"rows"`
					MaxSeats int      `json:"maxSeats"`
					Sides    []string `json:"sides"`
					Sectors  []string `json:"sectors"`
					Map      string   `json:"map"`
				} `json:"place"`
				RentalPeriodId int    `json:"rental_period_id"`
				RentalPeriod   string `json:"rental_period"`
				StatusName     string `json:"status_name"`
				Status         string `json:"status"`
			} `json:"event"`
			ShowId           interface{} `json:"show_id"`
			ShowName         string      `json:"show_name"`
			DateTime         string      `json:"date_time"`
			Amount           string      `json:"amount"`
			Status           string      `json:"status"`
			ReservedDateTime string      `json:"reserved_date_time"`
			StatusName       string      `json:"status_name"`
			Funds            []struct {
				Id        int       `json:"id"`
				Name      string    `json:"name"`
				Property  string    `json:"property"`
				Capacity  string    `json:"capacity"`
				CreatedAt time.Time `json:"created_at"`
				UpdatedAt time.Time `json:"updated_at"`
				Type      string    `json:"type"`
				Units     int       `json:"units"`
			} `json:"funds"`
			Note              interface{} `json:"note"`
			AvailableStatuses []struct {
				Id    string `json:"id"`
				Label string `json:"label"`
			} `json:"available_statuses"`
			Refund struct {
				Id            int         `json:"id"`
				OrderId       int         `json:"order_id"`
				SellerId      int         `json:"seller_id"`
				ClientId      interface{} `json:"client_id"`
				ContactId     interface{} `json:"contact_id"`
				EventId       interface{} `json:"event_id"`
				DateTime      string      `json:"date_time"`
				Amount        string      `json:"amount"`
				Reason        string      `json:"reason"`
				ReasonComment string      `json:"reason_comment"`
				Status        string      `json:"status"`
				CreatedAt     time.Time   `json:"created_at"`
				UpdatedAt     time.Time   `json:"updated_at"`
				DeletedAt     interface{} `json:"deleted_at"`
			} `json:"refund"`
			Token        string      `json:"token"`
			CountTickets interface{} `json:"count_tickets"`
			CashboxId    int         `json:"cashbox_id"`
			Cashboxes    []struct {
				Id          int         `json:"id"`
				Name        string      `json:"name"`
				SellerId    int         `json:"seller_id"`
				Address     string      `json:"address"`
				SellerName  string      `json:"seller_name"`
				CompanyId   interface{} `json:"company_id"`
				CompanyName string      `json:"company_name"`
			} `json:"cashboxes"`
			EventSeats []struct {
				Id         int         `json:"id"`
				SeatNumber int         `json:"seat_number"`
				RowSector  int         `json:"row_sector"`
				SectorId   int         `json:"sector_id"`
				SectorName string      `json:"sector_name"`
				Status     string      `json:"status"`
				SeatId     int         `json:"seat_id"`
				FundId     interface{} `json:"fund_id"`
				Price      int         `json:"price"`
				PriceZone  int         `json:"price_zone"`
				Seat       struct {
					Id         int    `json:"id"`
					HallId     int    `json:"hall_id"`
					HallName   string `json:"hall_name"`
					SectorId   int    `json:"sector_id"`
					SectorSlug string `json:"sector_slug"`
					Zona       string `json:"zona"`
					RowSector  int    `json:"row_sector"`
					SeatNumber int    `json:"seat_number"`
					StatusName string `json:"status_name"`
					CreatedAt  string `json:"created_at"`
				} `json:"seat"`
			} `json:"event_seats"`
			EventSeatsCount int          `json:"event_seats_count"`
			Tickets         []TicketData `json:"tickets"`
		} `json:"order"`
		Seller struct {
			Id          int    `json:"id"`
			ParentId    int    `json:"parent_id"`
			Name        string `json:"name"`
			Type        string `json:"type"`
			CompanyId   int    `json:"company_id"`
			ContactId   int    `json:"contact_id"`
			CompanyName string `json:"company_name"`
			ContactName string `json:"contact_name"`
			ParentName  string `json:"parent_name"`
		} `json:"seller"`
		Client  interface{}  `json:"client"`
		Contact interface{}  `json:"contact"`
		Event   interface{}  `json:"event"`
		Tickets []TicketData `json:"tickets"`
		Show    interface{}  `json:"show"`
	} `json:"data"`
}

type Sell struct {
	Data struct {
		Id       int         `json:"id"`
		OrderId  int         `json:"order_id"`
		SellerId int         `json:"seller_id"`
		AgentId  interface{} `json:"agent_id"`
		Seller   struct {
			Id        int         `json:"id"`
			ParentId  int         `json:"parent_id"`
			Name      string      `json:"name"`
			Type      string      `json:"type"`
			CompanyId int         `json:"company_id"`
			ContactId int         `json:"contact_id"`
			DeletedAt interface{} `json:"deleted_at"`
			CreatedAt interface{} `json:"created_at"`
			UpdatedAt time.Time   `json:"updated_at"`
			Contact   struct {
				Id           int         `json:"id"`
				DeletedAt    interface{} `json:"deleted_at"`
				CreatedAt    time.Time   `json:"created_at"`
				UpdatedAt    time.Time   `json:"updated_at"`
				Phone        string      `json:"phone"`
				Email        string      `json:"email"`
				Gender       string      `json:"gender"`
				DateBirth    interface{} `json:"date_birth"`
				PassNum      interface{} `json:"pass_num"`
				PassDate     interface{} `json:"pass_date"`
				PassIssued   interface{} `json:"pass_issued"`
				PassDivision interface{} `json:"pass_division"`
				Inn          int64       `json:"inn"`
				Snils        interface{} `json:"snils"`
				Address      interface{} `json:"address"`
				Type         string      `json:"type"`
				Name         string      `json:"name"`
				LastName     string      `json:"last_name"`
				MiddleName   string      `json:"middle_name"`
			} `json:"contact"`
		} `json:"seller"`
		SellerName  string `json:"seller_name"`
		SellerType  string `json:"seller_type"`
		ContactId   int    `json:"contact_id"`
		ContactName string `json:"contact_name"`
		ClientId    int    `json:"client_id"`
		ClientName  string `json:"client_name"`
		Type        string `json:"type"`
		PaymentType string `json:"payment_type"`
		User        struct {
			Id              int         `json:"id"`
			Name            string      `json:"name"`
			Email           string      `json:"email"`
			EmailVerifiedAt interface{} `json:"email_verified_at"`
			CreatedAt       time.Time   `json:"created_at"`
			UpdatedAt       time.Time   `json:"updated_at"`
			Position        interface{} `json:"position"`
			Organization    interface{} `json:"organization"`
			Avatar          interface{} `json:"avatar"`
			ContactId       int         `json:"contact_id"`
			SellerId        interface{} `json:"seller_id"`
			Timezone        string      `json:"timezone"`
			CashboxId       interface{} `json:"cashbox_id"`
		} `json:"user"`
		UserId  int `json:"user_id"`
		EventId int `json:"event_id"`
		Event   struct {
			Id            int    `json:"id"`
			Name          string `json:"name"`
			Datetime      string `json:"datetime"`
			Type          string `json:"type"`
			TimeLength    string `json:"timeLength"`
			AgeLimit      string `json:"ageLimit"`
			ShowId        int    `json:"show_id"`
			HallId        int    `json:"hall_id"`
			PlaceHallName string `json:"place_hall_name"`
			Promoter      struct {
				Name string `json:"name"`
			} `json:"promoter"`
			TicketFund     string `json:"ticketFund"`
			TicketProvider struct {
				Name string `json:"name"`
			} `json:"ticketProvider"`
			TicketBook int `json:"ticketBook"`
			Tickets    struct {
				Total int `json:"total"`
			} `json:"tickets"`
			Seats struct {
				Total int `json:"total"`
			} `json:"seats"`
			Place struct {
				Id       int      `json:"id"`
				Name     string   `json:"name"`
				Rows     int      `json:"rows"`
				MaxSeats int      `json:"maxSeats"`
				Sides    []string `json:"sides"`
				Sectors  []string `json:"sectors"`
				Map      string   `json:"map"`
			} `json:"place"`
			RentalPeriodId int    `json:"rental_period_id"`
			RentalPeriod   string `json:"rental_period"`
			StatusName     string `json:"status_name"`
			Status         string `json:"status"`
		} `json:"event"`
		ShowId           int    `json:"show_id"`
		ShowName         string `json:"show_name"`
		DateTime         string `json:"date_time"`
		Amount           string `json:"amount"`
		Status           string `json:"status"`
		ReservedDateTime string `json:"reserved_date_time"`
		StatusName       string `json:"status_name"`
		Funds            []struct {
			Id        int        `json:"id"`
			Name      string     `json:"name"`
			Property  string     `json:"property"`
			Capacity  string     `json:"capacity"`
			CreatedAt *time.Time `json:"created_at"`
			UpdatedAt *time.Time `json:"updated_at"`
			Type      string     `json:"type"`
			Units     int        `json:"units"`
		} `json:"funds"`
		Note              interface{} `json:"note"`
		AvailableStatuses []struct {
			Id    string `json:"id"`
			Label string `json:"label"`
		} `json:"available_statuses"`
		Refund       interface{} `json:"refund"`
		Token        string      `json:"token"`
		CountTickets interface{} `json:"count_tickets"`
		CashboxId    interface{} `json:"cashbox_id"`
		Cashboxes    []struct {
			Id          int         `json:"id"`
			Name        string      `json:"name"`
			SellerId    int         `json:"seller_id"`
			Address     string      `json:"address"`
			SellerName  string      `json:"seller_name"`
			CompanyId   interface{} `json:"company_id"`
			CompanyName string      `json:"company_name"`
		} `json:"cashboxes"`
		EventSeats []struct {
			Id         int         `json:"id"`
			SeatNumber int         `json:"seat_number"`
			RowSector  int         `json:"row_sector"`
			SectorId   int         `json:"sector_id"`
			SectorName string      `json:"sector_name"`
			Status     string      `json:"status"`
			SeatId     int         `json:"seat_id"`
			FundId     interface{} `json:"fund_id"`
			Price      int         `json:"price"`
			PriceZone  int         `json:"price_zone"`
			Seat       struct {
				Id         int    `json:"id"`
				HallId     int    `json:"hall_id"`
				HallName   string `json:"hall_name"`
				SectorId   int    `json:"sector_id"`
				SectorSlug string `json:"sector_slug"`
				Zona       string `json:"zona"`
				RowSector  int    `json:"row_sector"`
				SeatNumber int    `json:"seat_number"`
				StatusName string `json:"status_name"`
				CreatedAt  string `json:"created_at"`
			} `json:"seat"`
		} `json:"event_seats"`
		EventSeatsCount int          `json:"event_seats_count"`
		Tickets         []TicketData `json:"tickets"`
	} `json:"data"`
}

type TicketData struct {
	Id          int         `json:"id"`
	Number      string      `json:"number"`
	Status      string      `json:"status"`
	StatusName  string      `json:"status_name"`
	Amount      int         `json:"amount"`
	OrderId     int         `json:"order_id"`
	ShowId      int         `json:"show_id"`
	Show        string      `json:"show"`
	CashboxId   interface{} `json:"cashbox_id"`
	CashboxName string      `json:"cashbox_name"`
	EventSeatId int         `json:"event_seat_id"`
	SeatNumber  int         `json:"seat_number"`
	RowSector   int         `json:"row_sector"`
	Zona        string      `json:"zona"`
	Event       struct {
		Id             int       `json:"id"`
		HallId         int       `json:"hall_id"`
		ShowId         int       `json:"show_id"`
		DateTime       string    `json:"date_time"`
		CreatedAt      time.Time `json:"created_at"`
		UpdatedAt      time.Time `json:"updated_at"`
		Status         string    `json:"status"`
		RentalPeriodId int       `json:"rental_period_id"`
		Show           struct {
			Id           int       `json:"id"`
			Type         string    `json:"type"`
			Name         string    `json:"name"`
			Description  string    `json:"description"`
			Duration     string    `json:"duration"`
			Intermission bool      `json:"intermission"`
			AgeLimit     string    `json:"age_limit"`
			CreatedAt    time.Time `json:"created_at"`
			UpdatedAt    time.Time `json:"updated_at"`
			ShowOnIndex  bool      `json:"show_on_index"`
			Poster       string    `json:"poster"`
			ProgramId    int       `json:"program_id"`
		} `json:"show"`
	} `json:"event"`
	ScannedAt    interface{} `json:"scanned_at"`
	ScanCount    int         `json:"scan_count"`
	UserId       int         `json:"user_id"`
	UserFullname string      `json:"user_fullname"`
	Date         interface{} `json:"date"`
}
