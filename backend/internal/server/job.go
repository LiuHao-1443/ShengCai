package server

import (
	"context"
	"shengcai/internal/repository"
	"shengcai/internal/service"
	"shengcai/pkg/log"
	"time"
)

type Job struct {
	log             *log.Logger
	shengcaiService service.ShengCaiService
	aiRepository    repository.AIRepository
}

func NewJob(
	log *log.Logger,
	shengcaiService service.ShengCaiService,
	aiRepository repository.AIRepository,
) *Job {
	return &Job{
		log:             log,
		shengcaiService: shengcaiService,
		aiRepository:    aiRepository,
	}
}
func (j *Job) Start(ctx context.Context) error {
	for {
		// 调用 CreateData 方法
		if err := j.shengcaiService.CreateData(ctx); err != nil {
			return err // 如果需要处理错误，可以在这里处理
		}

		// 等待 10 秒
		time.Sleep(5 * 60 * time.Second)
	}

	//abstract, keyword, _ := j.aiRepository.GenerateAbstractAndKeyword(ctx, "上周很火的“酱香拿铁” 大家都喝了吗？\n上周很火的“酱香拿铁” 大家都喝了吗？\n\n我在朋友圈看到不少圈友已经在研究它为什么那么火，为什么赚钱了，而像这种爆火的品牌有很多，像茶百道、蜜雪冰城、暴打柠檬茶、霸王茶姬等等，都是怎么火起来、怎么赚到钱的呢？\n\n不久前，圈友@易生 和大家分享了一篇文章：《线下实体复苏，经营一家奶茶店到底挣不挣钱？》，受到很多小伙伴的关注，评论区也有很多疑问：\n\n应季柠檬茶冬天营收会下降几成？\n被困在自营奶茶店上，生意惨淡想关闭不知道如何破局\n一家蜜雪冰城加盟一年大概赚多少？\n有没有选址的方法？\n小白是否建议加盟做开奶茶店？\n......\n\n所以这周三，我们也邀请了@易生 作为实战家，与生财实战访谈官@亮哥 一起深入聊一聊，线下实体生意复苏，开一家奶茶店、餐饮店到底挣不挣钱。\n\n两位都是有线下餐饮实体店经验的实战家，易生是经营奶茶店，而亮哥是经营火锅店，他们对于线下餐饮实体怎么选址、品牌怎么选，如何运营等等都有不少的经验分享，帮助大家避坑。\n\n如果你也想创业开一家实体餐饮店，那周三（13日）晚 20：00，生财有术直播间，一定要来。\n展开全部\nVvoPbBTX5os00pxm4bpckmoCn0e.jpg\n\n查看详情\n维倪、黄小鱼?、雅俊、查克、张有财、艾小飞、玟、蒋儒钢、岁月静好、薯条\n 等29人觉得很赞\n鱼丸 | 亦仁助理：置顶大家有想要了解的内容，或者有疑问的地方，欢迎评论区留言，直播时一一为大家解答呀\nAAk3b0fobopCOOxN6GJcjcd1nvb.png\n\n2023-09-11 14:57\nMazc：喝了，我觉得很不好喝\nIc2Db96j6oxykFxI1Eych8SVnlg.png\n\n2023-09-11 14:13\n小满：期待住了\n2023-09-11 14:23\n倾欣为红颜：太急时了，我正在发愁怎么提升营业额。加盟的某品牌，一杯均价15元，从四月到现在，堂食+外卖营业额共8w+，进了9月更是断崖式下跌，一天卖200来块。\n2023-09-11 14:31\n")
	//fmt.Println("abstract ==>", abstract)
	//fmt.Println("keyword ==>", keyword)

	// 这个代码不会被执行到，因为前面的 for 是一个死循环
	return nil
}
func (j *Job) Stop(ctx context.Context) error {
	return nil
}
