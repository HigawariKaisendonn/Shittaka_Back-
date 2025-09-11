package services

import (
	"context"
	"log"
	"os"
	"testing"

	entities "Shittaka_back/internal/domain/choices/entities"
	"Shittaka_back/internal/domain/choices/repositories"

	"github.com/joho/godotenv"
	"github.com/nedpals/supabase-go"
	"github.com/stretchr/testify/assert"
)

func TestChoiceService_CRUD(t *testing.T) {
	// .envを読み込む
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on system environment variables")
	}

	url := os.Getenv("SUPABASE_URL")
	serviceRoleKey := os.Getenv("SUPABASE_SERVICE_ROLE_KEY")
	if url == "" || serviceRoleKey == "" {
		t.Fatal("SUPABASE_URL or SUPABASE_SERVICE_ROLE_KEY is not set")
	}

	// Supabase クライアント作成
	client := supabase.CreateClient(url, serviceRoleKey)

	// リポジトリとサービス作成
	repo := repositories.NewChoiceRepository(client)
	service := NewChoiceService(repo)

	ctx := context.Background()

	// ---------------------------
	// 1. Create
	// ---------------------------
	newChoice := entities.Choice{
		QuestionID: 1, // 実在する question_id を使用
		Text:       "テスト選択肢",
		IsCorrect:  false,
	}

	created, err := service.CreateChoice(ctx, newChoice)
	assert.NoError(t, err)
	assert.NotNil(t, created)
	t.Logf("Created Choice: %+v", created)

	// ---------------------------
	// 2. GetByQuestionID
	// ---------------------------
	choices, err := service.GetChoices(ctx, int64(newChoice.QuestionID))
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(choices), 1)
	t.Logf("Choices for QuestionID=%d: %+v", newChoice.QuestionID, choices)

	// ---------------------------
	// 3. Update
	// ---------------------------
	created.Text = "更新済み選択肢"
	updated, err := service.UpdateChoice(ctx, *created)
	assert.NoError(t, err)
	assert.Equal(t, "更新済み選択肢", updated.Text)
	t.Logf("Updated Choice: %+v", updated)

	// ---------------------------
	// 4. Delete
	// ---------------------------
	err = service.DeleteChoice(ctx, int64(updated.ID))
	assert.NoError(t, err)
	t.Logf("Deleted Choice ID: %d", updated.ID)
}
