"use client";

import { Card, CardContent, CardHeader } from '@/components/ui/card';
import { PlaceholdersAndVanishInput } from '@/components/ui/placeholders-and-vanish-input'
import React, { useState } from 'react'
import axios from "axios"

export default function AIChat() {
  const [input, setInput] = useState("");
  const [loading, setLoading] = useState(false);
  const [aIResponse, setAIResponse] = useState("");
  const placeholders = [
    "What's the first rule of Fight Club?",
    "Who is Tyler Durden?",
    "Where is Andrew Laeddis Hiding?",
    "Write a Javascript method to reverse a string",
    "How to assemble your own PC?",
  ];

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setInput(e.target.value);
  };
  const onSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setLoading(true);
    const response = await axios.post(`${process.env.NEXT_PUBLIC_SERVER_URL}/generate`,{model: 'llama2-7b',prompt: input},{headers: {
      'Authorization': `Bearer ${process.env.NEXT_PUBLIC_TEST_TOKEN}`}})
    setAIResponse(response.data.data.output);
    setLoading(false);
  };
  return (
     <div className="h-[40rem] flex flex-col justify-center  items-center px-4">
      <h2 className="mb-10 sm:mb-20 text-xl text-center sm:text-5xl dark:text-white text-dark">
        Ask Aceternity UI Anything
      </h2>
      <PlaceholdersAndVanishInput
        placeholders={placeholders}
        onChange={handleChange}
        onSubmit={onSubmit}
      />
      {loading ? (
        <p className="text-center text-2xl">Loading...</p>
      ) : (
        <Card className="mt-10">
          <CardHeader>Response</CardHeader>
          <CardContent>{aIResponse}</CardContent>
        </Card>
      )}
    </div> 
  )
}

