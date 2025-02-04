"use client";

import { Card, CardContent, CardDescription } from "../../components/ui/card"; // Updated import path
import React, { useState } from "react";
import axios from "axios";
import Image from "next/image";
import { AIInputWithLoading } from "@/components/ui/ai-input-with-loading";

type AIResponse = {
  requestID: string;
  aIResponse: string;
  imageURL: string;
};
export default function AIChat() {
  const [responses, setResponses] = useState<AIResponse[]>([]);

  const onSubmit = async (input: string) => {
    const token = "";
    const response = await axios.post(
      `${process.env.NEXT_PUBLIC_SERVER_URL}/generate`,
      { model: "llama2-7b", prompt: input },
      {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      },
    );
    setResponses((prev) => [
      {
        requestID: response.data.data.request_id as string,
        aIResponse: response.data.data.output as string,
        imageURL: response.data.data.image_url as string,
      },
      ...prev,
    ]);
  };

  return (
    <div className="min-h-screen flex flex-col justify-center items-center px-4 bg-zinc-800">
      <h2 className="mb-10 sm:mb-20 text-5xl text-center sm:text-5xl dark:text-white text-dark font-bold">
        Ask me to generate a social media post for you!
      </h2>
      <AIInputWithLoading
        onSubmit={onSubmit}
        placeholder="Ask me anything..."
        autoAnimate={false}
      />
      <div className="flex flex-col gap-10 justify-center items-center">
        {responses.map((response) => {
          return (
            <CardWithImage
              key={response.requestID}
              imageURL={response.imageURL}
              aIResponse={response.aIResponse}
            />
          );
        })}
      </div>
    </div>
  );
}

type CardWithImageProps = {
  imageURL: string;
  aIResponse: string;
};
const CardWithImage = ({ imageURL, aIResponse }: CardWithImageProps) => {
  return (
    <Card className="flex flex-col justify-center items-center p-10 w-1/2 bg-zinc-900">
      <CardContent>
        <Image
          src={imageURL}
          className="rounded-xl shadow-lg"
          alt="AI Image"
          width={500}
          height={500}
        />
      </CardContent>
      <CardDescription className="text-md font-semibold text-white border-t-4 border-zinc-700 pt-10">
        {aIResponse}
      </CardDescription>
    </Card>
  );
};
