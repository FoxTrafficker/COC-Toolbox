"use client";

import React, {useState} from "react";
import {Box, Paper} from "@mui/material";


function BattleIndicator() {
    let l = [1, 2, 3, 4, 5, 6, 7, 8]
    let tar = [2, 5, 7, 10]


    return (
        <Box>
            <Paper>{l.filter(num => tar.includes(num))}</Paper>
        </Box>
    )
}

export default BattleIndicator;