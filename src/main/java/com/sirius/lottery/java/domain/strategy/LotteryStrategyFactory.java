package com.sirius.lottery.java.domain.strategy;

import org.springframework.stereotype.Component;

import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;

@Component
public class LotteryStrategyFactory {

    private final Map<String, LotteryStrategy> strategyMap;

    public LotteryStrategyFactory(Map<String, LotteryStrategy> strategyMap) {
        this.strategyMap = new ConcurrentHashMap<>(strategyMap);
    }

    public LotteryStrategy getStrategy(String strategyName) {
        return strategyMap.get(strategyName);
    }
}
