package service

import (
	"context"
	"fms_audit/internal/domain"
	"fms_audit/internal/infrastructure/repository"
)

type DishService struct {
	dishRepo *repository.DishRepository
}

func NewDishService(dishRepo *repository.DishRepository) *DishService {
	return &DishService{
		dishRepo: dishRepo,
	}
}

func (s *DishService) GetDishes(ctx context.Context, filter domain.DishFilter) (domain.DishListResponse, error) {
	dishes, total, err := s.dishRepo.GetAll(ctx, filter)
	if err != nil {
		return domain.DishListResponse{}, err
	}

	totalPages := repository.CalculateTotalPages(total, filter.Limit)

	response := domain.DishListResponse{
		Success: true,
	}
	response.Data.Dishes = dishes
	response.Data.Pagination = domain.PaginationResponse{
		Page:       filter.Page,
		Limit:      filter.Limit,
		Total:      total,
		TotalPages: totalPages,
	}

	return response, nil
}

func (s *DishService) GetFeaturedDishes(ctx context.Context) ([]domain.DishCardResponse, error) {
	dishes, err := s.dishRepo.GetFeatured(ctx, 2)
	if err != nil {
		return nil, err
	}

	var cardResponses []domain.DishCardResponse
	for _, dish := range dishes {
		cardResponses = append(cardResponses, s.convertToCardResponse(dish, true))
	}

	return cardResponses, nil
}

func (s *DishService) GetDishesForFrontend(ctx context.Context, filter domain.DishFilter) (domain.DishCardListResponse, error) {
	dishes, total, err := s.dishRepo.GetAll(ctx, filter)
	if err != nil {
		return domain.DishCardListResponse{}, err
	}

	totalPages := repository.CalculateTotalPages(total, filter.Limit)

	// Convert to DishCardResponse
	var cardResponses []domain.DishCardResponse
	for _, dish := range dishes {
		cardResponses = append(cardResponses, s.convertToCardResponse(dish, false))
	}

	response := domain.DishCardListResponse{
		Success: true,
	}
	response.Data.Dishes = cardResponses
	response.Data.Pagination = domain.PaginationResponse{
		Page:       filter.Page,
		Limit:      filter.Limit,
		Total:      total,
		TotalPages: totalPages,
	}

	return response, nil
}

func (s *DishService) GetDishesByStirFried(ctx context.Context) ([]domain.DishCardResponse, error) {
	dishes, err := s.dishRepo.GetStirFriedDishes(ctx)
	if err != nil {
		return nil, err
	}

	var cardResponses []domain.DishCardResponse
	for _, dish := range dishes {
		cardResponses = append(cardResponses, s.convertToCardResponse(dish, false))
	}

	return cardResponses, nil
}

func (s *DishService) GetDishesBySteamed(ctx context.Context) ([]domain.DishCardResponse, error) {
	dishes, err := s.dishRepo.GetSteamedDishes(ctx)
	if err != nil {
		return nil, err
	}

	var cardResponses []domain.DishCardResponse
	for _, dish := range dishes {
		cardResponses = append(cardResponses, s.convertToCardResponse(dish, false))
	}

	return cardResponses, nil
}

func (s *DishService) GetDishesByGrilled(ctx context.Context) ([]domain.DishCardResponse, error) {
	dishes, err := s.dishRepo.GetGrilledDishes(ctx)
	if err != nil {
		return nil, err
	}

	var cardResponses []domain.DishCardResponse
	for _, dish := range dishes {
		cardResponses = append(cardResponses, s.convertToCardResponse(dish, false))
	}

	return cardResponses, nil
}

func (s *DishService) GetDrinks(ctx context.Context) ([]domain.DishCardResponse, error) {
	dishes, err := s.dishRepo.GetDrinks(ctx)
	if err != nil {
		return nil, err
	}

	var cardResponses []domain.DishCardResponse
	for _, dish := range dishes {
		cardResponses = append(cardResponses, s.convertToCardResponse(dish, false))
	}

	return cardResponses, nil
}

func (s *DishService) convertToCardResponse(dish domain.Dish, showDetails bool) domain.DishCardResponse {
	return domain.DishCardResponse{
		ID:            dish.ID,
		Img:           dish.ImageURL,
		Title:         dish.Name,
		Rating:        &dish.Rating,
		Price:         &dish.Price,
		Description:   &dish.Description,
		CookingMethod: &dish.CookingMethod,
		ShowDetails:   &showDetails,
	}
}
