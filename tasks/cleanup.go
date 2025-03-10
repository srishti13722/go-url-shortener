package tasks

import (
	"context"
	"fmt"
	"go-url-shortener/database"
	"time"
)

func DeleteExpiredUrls(){
	for{
		//delete urls data have expired
		_ , err := database.DB.Exec(context.Background(), "DELETE FROM urls WHERE expiry < NOW()")
		if err != nil {
			fmt.Println("Failed to delete expired URLs:", err)
		} else {
			fmt.Println("Expired URLs deleted successfully")
		}

		// Sleep for 1 hour before running again
		time.Sleep(1 * time.Hour)
	}
}