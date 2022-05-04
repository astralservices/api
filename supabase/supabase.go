package db

import (
	"os"

	"github.com/nedpals/supabase-go"
)

func New() *supabase.Client {
	supabaseClient := supabase.CreateClient(os.Getenv("SUPABASE_URL"), os.Getenv("SUPABASE_KEY"))
	return supabaseClient
}
