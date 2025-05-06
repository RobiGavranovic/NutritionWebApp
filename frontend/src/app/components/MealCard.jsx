import Card from "@mui/material/Card";
import CardContent from "@mui/material/CardContent";
import CardMedia from "@mui/material/CardMedia";
import CardActionArea from "@mui/material/CardActionArea";
import Typography from "@mui/material/Typography";

export default function MealCard({ meal }) {
  return (
    <Card
      key={meal.IDMeal}
      sx={{
        maxWidth: 300,
        height: { md: 444 },
        display: "flex",
        flexDirection: "column",
      }}
    >
      <CardActionArea
        sx={{
          flexGrow: 1,
          display: "flex",
          flexDirection: "column",
          alignItems: "stretch",
        }}
      >
        <CardMedia
          component="img"
          height="180"
          image={meal.StrMealThumb}
          alt={meal.StrMeal}
        />
        <CardContent
          sx={{ flexGrow: 1, display: "flex", flexDirection: "column" }}
        >
          <Typography gutterBottom variant="h6" component="div" align="center">
            {meal.StrMeal}
          </Typography>
          <Typography
            variant="body2"
            sx={{
              color: "text.secondary",
              overflow: "hidden",
              display: "-webkit-box",
              WebkitLineClamp: 4,
              WebkitBoxOrient: "vertical",
            }}
          >
            {meal.StrInstructions}
          </Typography>
        </CardContent>
      </CardActionArea>
    </Card>
  );
}
