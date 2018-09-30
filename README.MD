# Booking Ticket project

## Setup env for dev
```
docker-compose up -d 
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
