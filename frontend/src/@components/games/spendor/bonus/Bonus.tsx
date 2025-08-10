import CardContent from "@mui/material/CardContent"
import Card from "@mui/material/Card"
import Box from "@mui/material/Box"

import type { Bonus as SplendorBonus } from "@types/splendor/games"

export type BonusProps = {
  bonus: SplendorBonus
}

const Bonus = (props: BonusProps) => {
  const image = "/images/splendor/bonuses/" + props.bonus.id + ".png"

  const imageFormating = {
    height: 60,
    width: 60,
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
            height: 60,
            width: 60,
            maxHeight: { xs: 60, md: 60 },
            maxWidth: { xs: 60, md: 60 },
          }}
          src={image}
        />
      </CardContent>
    </Card>
  )
}

export default Bonus 
