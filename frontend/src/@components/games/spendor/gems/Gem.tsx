import { GemCount } from "@/@types/splendor/games"


export type GemProps = {
  gem: GemCount
}

const Gem = (props: GemProps) => {
  const { gem } = props

  return <div>Gem</div>
}

export default Gem
