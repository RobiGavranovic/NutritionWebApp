
import { motion } from 'framer-motion';

import Card from '@mui/material/Card';
import CardContent from '@mui/material/CardContent';
import CardMedia from '@mui/material/CardMedia';
import CardActionArea from '@mui/material/CardActionArea';
import Typography from '@mui/material/Typography';

export default function MealCard({ meal }) {
    return (




        <Card key={meal.IDMeal} sx={{ maxWidth: 300 }}>
            <CardActionArea>
            <CardMedia
                component="img"
                height="180"
                image={meal.StrMealThumb}
                alt={meal.StrMeal}
            />
            <CardContent>
                <Typography gutterBottom variant="h5" component="div" align="center">
                {meal.StrMeal}
                </Typography>
                <Typography variant="body2" sx={{ color: 'text.secondary' }}>
                {meal.StrInstructions.slice(0, 100)}...
                </Typography>
            </CardContent>
            </CardActionArea>
        </Card>




    );
}
