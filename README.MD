# Booking Ticket project

## Setup env for dev
```
docker-compose up -d 
go run intalTicket/cmd/main.go // run this for inital ticket for testing
```

## Format Data 
```
current_round : string
 -- i.e. "1"
r_{round}:{Seat} 
 -- The available seat .
r_{round}_u:{Seat} 
 -- The booked seat but didn't confirm yet.

```

## Main function
- POST: Book the seat  
  - PATH
    - /book
  - Parameter
    - seat: string (seat no.)
  - Return 
    - status 200
    - Field
      - success: boolean
      - round: Number (Round number of this ticket)
      - seat: string (seat of this ticket)
      - reserve_expired_time: number (Reservation expired time in timestamp)
- POST: Confirm Ticket
  - PATH
    - /confirm
  - Parameter
    - seat
  - Return
    - status 200
    
- POST: Cancel Ticket /cancel
  - PATH
    - /cancel
  - Parameter
    - seat: string (seat no.)
  - Return
    - status 200 or 403
    - Field
      - success: boolean
- GET: Get Remaining Seat(s)
  - PATH
    - /remaining
  - Return
    - status 200
    - Field
      - seats: string[] 
      - unconfimedTicketsCount: Number
- GET: Get all booked tickets
  - PATH
    - /bookings
  - Return
    - status 200
    - Field
      - bookings: string[]  (list of all booking)


## TODO
- Inital ticket  {DONE}
- Get remain ticket {DONE}
- Book seat 
- Confirm Seat
- Cancel Seat
- Add test for all


## Use case
- 1. If all ticket booked but didn't confirm yet not allow to book in the next round
