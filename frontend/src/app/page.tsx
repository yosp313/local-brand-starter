import { PricingCards } from "@/components/pricing-cards";
import { Feature } from "@/components/ui/feature-with-advantages";
import { Footerdemo } from "@/components/ui/footer-section";
import { HeroGeometric } from "@/components/ui/shape-landing-hero";

export default function Home() {
  const tiers = [
    {
      name: "Free",
      price: 0,
      description: "For small businesses that are just getting started.",
      features: [
        { name: "Unlimited products", included: true },
        { name: "Unlimited orders", included: false },
        { name: "24/7 support", included: false },
        { name: "Custom domain", included: false },
      ],
    },
    {
      name: "Basic",
      price: 29,
      description: "For businesses that are looking to grow.",
      features: [
        { name: "Unlimited products", included: true },
        { name: "Unlimited orders", included: true },
        { name: "24/7 support", included: true },
        { name: "Custom domain", included: false },
      ],
    },
    {
      name: "Pro",
      price: 49,
      description: "For businesses that are looking to scale.",
      features: [
        { name: "Unlimited products", included: true },
        { name: "Unlimited orders", included: true },
        { name: "24/7 support", included: true },
        { name: "Custom domain", included: true },
      ],
    },
  ];

  return (
    <>
      <HeroGeometric
        badge="Local Brand Starter"
        title1="Grow Your"
        title2="Brand With Us"
      />
      <Feature />
      <h1 className="text-4xl font-bold text-center mt-12">Pricing</h1>
      <PricingCards tiers={tiers} className="gap-6" />
      <Footerdemo />
    </>
  );
}
