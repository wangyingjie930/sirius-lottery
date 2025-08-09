package com.sirius.lottery.java.application.service;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.sirius.lottery.java.application.dto.DrawRequest;
import com.sirius.lottery.java.application.dto.DrawResponse;
import com.sirius.lottery.java.application.dto.LotteryInstanceResponse;
import com.sirius.lottery.java.application.dto.StockActionRequest;
import com.sirius.lottery.java.domain.entity.Instance;
import com.sirius.lottery.java.domain.entity.Pool;
import com.sirius.lottery.java.domain.entity.Prize;
import com.sirius.lottery.java.domain.entity.WinRecord;
import com.sirius.lottery.java.domain.model.DrawContext;
import com.sirius.lottery.java.domain.repository.InstanceRepository;
import com.sirius.lottery.java.domain.repository.PoolRepository;
import com.sirius.lottery.java.domain.repository.PrizeRepository;
import com.sirius.lottery.java.domain.repository.WinRecordRepository;
import com.sirius.lottery.java.domain.strategy.LotteryStrategy;
import com.sirius.lottery.java.domain.strategy.LotteryStrategyFactory;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.data.redis.core.StringRedisTemplate;
import org.springframework.data.redis.core.script.DefaultRedisScript;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.time.LocalDateTime;
import java.util.Collections;
import java.util.List;
import java.util.UUID;
import java.util.concurrent.TimeUnit;
import java.util.stream.Collectors;

@Slf4j
@Service
@RequiredArgsConstructor
public class LotteryServiceImpl implements LotteryService {

    private static final String DEDUCT_STOCK_LUA_SCRIPT =
            "local stock = redis.call('GET', KEYS[1])\n" +
            "if not stock or tonumber(stock) < tonumber(ARGV[1]) then\n" +
            "    return 0\n" +
            "end\n" +
            "redis.call('DECRBY', KEYS[1], ARGV[1])\n" +
            "return 1";

    private final InstanceRepository instanceRepository;
    private final PoolRepository poolRepository;
    private final PrizeRepository prizeRepository;
    private final WinRecordRepository winRecordRepository;
    private final LotteryStrategyFactory strategyFactory;
    private final StringRedisTemplate redisTemplate;
    private final ObjectMapper objectMapper;

    @Override
    @Transactional
    public DrawResponse draw(DrawRequest req) {
        long userId = 100L;

        Instance instance = getAndCacheInstance(req.getInstanceId());

        if (instance.getStatus() != 2) { // 2: Active
            throw new RuntimeException("Lottery is not active");
        }
        LocalDateTime now = LocalDateTime.now();
        if (now.isBefore(instance.getStartTime()) || now.isAfter(instance.getEndTime())) {
            throw new RuntimeException("Lottery is not within the active time range");
        }

        Pool pool = instance.getPools().stream().findFirst()
                .orElseThrow(() -> new RuntimeException("Lottery pool not found"));

        LotteryStrategy strategy = strategyFactory.getStrategy(pool.getLotteryStrategy());
        if (strategy == null) {
            throw new RuntimeException("Lottery strategy not found: " + pool.getLotteryStrategy());
        }

        DrawContext drawContext = new DrawContext(instance.getInstanceId(), userId, pool, pool.getPrizes());
        Prize wonPrize = strategy.draw(drawContext);

        if (wonPrize == null || (wonPrize.getIsSpecial() != null && wonPrize.getIsSpecial())) {
            return new DrawResponse("THANK_YOU_ORDER", wonPrize != null ? wonPrize.getPrizeId() : "THANK_YOU", false);
        }

        boolean stockDeducted = deductStockInRedis(instance.getInstanceId(), wonPrize.getPrizeId(), 1);
        if (!stockDeducted) {
            throw new RuntimeException("Prize stock is insufficient");
        }

        WinRecord record = createWinRecord(req, userId, wonPrize);

        return new DrawResponse(record.getOrderId(), record.getPrizeId(), true);
    }

    @Override
    public LotteryInstanceResponse getLotteryInstance(String instanceId) {
        Instance instance = getAndCacheInstance(instanceId);
        return mapToLotteryInstanceResponse(instance);
    }

    private Instance getAndCacheInstance(String instanceId) {
        String key = "lottery:instance:" + instanceId;
        try {
            String cachedInstance = redisTemplate.opsForValue().get(key);
            if (cachedInstance != null) {
                return objectMapper.readValue(cachedInstance, Instance.class);
            }
        } catch (Exception e) {
            log.error("Error reading instance from Redis cache", e);
        }

        Instance instance = instanceRepository.findByInstanceId(instanceId)
                .orElseThrow(() -> new RuntimeException("Lottery instance not found"));

        List<Pool> pools = poolRepository.findByInstanceId(instanceId);
        for (Pool pool : pools) {
            List<Prize> prizes = prizeRepository.findByPoolId(pool.getId());
            pool.setPrizes(prizes);
            // Cache stock for each prize
            for (Prize prize : prizes) {
                String stockKey = "lottery:stock:" + instanceId + ":" + prize.getPrizeId();
                redisTemplate.opsForValue().set(stockKey, prize.getAllocatedStock().toString());
            }
        }
        instance.setPools(pools);

        try {
            String instanceJson = objectMapper.writeValueAsString(instance);
            redisTemplate.opsForValue().set(key, instanceJson, 1, TimeUnit.HOURS);
        } catch (Exception e) {
            log.error("Error writing instance to Redis cache", e);
        }

        return instance;
    }

    private boolean deductStockInRedis(String instanceId, String prizeId, int num) {
        String key = "lottery:stock:" + instanceId + ":" + prizeId;
        DefaultRedisScript<Long> redisScript = new DefaultRedisScript<>(DEDUCT_STOCK_LUA_SCRIPT, Long.class);
        Long result = redisTemplate.execute(redisScript, Collections.singletonList(key), String.valueOf(num));
        return result != null && result == 1;
    }

    private WinRecord createWinRecord(DrawRequest req, long userId, Prize prize) {
        WinRecord record = WinRecord.builder()
                .orderId(UUID.randomUUID().toString())
                .requestId(req.getRequestId())
                .instanceId(req.getInstanceId())
                .userId(userId)
                .prizeId(prize.getPrizeId())
                .status(1) // 1: Pending
                .build();
        return winRecordRepository.save(record);
    }

    private LotteryInstanceResponse mapToLotteryInstanceResponse(Instance instance) {
        LotteryInstanceResponse response = new LotteryInstanceResponse();
        response.setInstanceId(instance.getInstanceId());
        response.setName(instance.getInstanceName());
        response.setStartTime(instance.getStartTime());
        response.setEndTime(instance.getEndTime());
        response.setServerTime(LocalDateTime.now());
        response.setTemplateStyle("default");
        response.setTemplateConfig(new LotteryInstanceResponse.TemplateConfig());

        response.setPools(instance.getPools().stream().map(p -> {
            LotteryInstanceResponse.Pool poolDto = new LotteryInstanceResponse.Pool();
            poolDto.setPoolId(p.getId().toString());
            poolDto.setPoolName(p.getPoolName());
            poolDto.setCost(Collections.emptyList());
            poolDto.setPrizes(p.getPrizes().stream().map(prize -> {
                LotteryInstanceResponse.Prize prizeDto = new LotteryInstanceResponse.Prize();
                prizeDto.setPrizeId(prize.getPrizeId());
                prizeDto.setPosition(0);
                return prizeDto;
            }).collect(Collectors.toList()));
            return poolDto;
        }).collect(Collectors.toList()));

        return response;
    }

    @Override
    public void deductStock(StockActionRequest req) {
        if (!deductStockInRedis(req.getInstanceId(), req.getPrizeId(), req.getNum())) {
            throw new RuntimeException("Insufficient stock");
        }
    }

    @Override
    public void increaseStock(StockActionRequest req) {
        String key = "lottery:stock:" + req.getInstanceId() + ":" + req.getPrizeId();
        redisTemplate.opsForValue().increment(key, req.getNum());
    }
}
