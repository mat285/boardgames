import CardContent from "@mui/material/CardContent"
import Card from "@mui/material/Card"
import Box from "@mui/material/Box"

import type { Card as SplendorCard } from "@/types/splendor/games"

export type CardProps = {
    card: SplendorCard
}

const SplendorCard = (props: CardProps) => {
  const image = "images/splendor/cards/" + props.card.id+".png"

  const imageFormating = {
    height: 130,
    width: 90,
  }


return (
    <Card sx={{
      padding: 0, 
      ...imageFormating,
    }}>
      <CardContent sx={{
      padding: 0, 
      ...imageFormating,
    }}>
      <Box
      component="img"
      sx={{
        height: 130,
        width: 90,
        maxHeight: { xs: 130, md: 130 },
        maxWidth: { xs: 90, md: 90 },
      }}
      src={image}
    />
      </CardContent>
    </Card>
  )
}

export default SplendorCard 
