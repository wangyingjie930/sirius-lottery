package com.sirius.lottery.java.domain.strategy;

import com.sirius.lottery.java.domain.entity.Prize;
import com.sirius.lottery.java.domain.model.DrawContext;
import org.springframework.stereotype.Component;

import java.util.Collections;
import java.util.Comparator;
import java.util.List;
import java.util.Random;

@Component("independent")
public class IndependentLotteryStrategy implements LotteryStrategy {

    private final Random random = new Random();

    @Override
    public Prize draw(DrawContext drawContext) {
        List<Prize> prizes = drawContext.getPrizes();
        if (prizes == null || prizes.isEmpty()) {
            return null;
        }

        // Sort prizes by probability, descending. While not strictly necessary for this algorithm's
        // correctness, it can be a good practice for some weighted selection algorithms.
        // We'll keep it for consistency with the original Go code's comment.
        prizes.sort(Comparator.comparing(Prize::getProbability).reversed());

        double randomValue = random.nextDouble();
        double cumulativeProbability = 0.0;

        for (Prize prize : prizes) {
            if (prize.getIsSpecial() != null && prize.getIsSpecial()) {
                continue; // "Thank you" prizes are not drawn by probability
            }
            cumulativeProbability += prize.getProbability();
            if (randomValue < cumulativeProbability) {
                return prize; // A prize is won
            }
        }

        // If no prize is won, return a "Thank You" prize if available
        return prizes.stream()
                .filter(p -> p.getIsSpecial() != null && p.getIsSpecial())
                .findFirst()
                .orElse(null); // Or null if no "Thank You" prize is configured
    }
}
